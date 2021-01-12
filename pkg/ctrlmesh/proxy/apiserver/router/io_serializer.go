/*
Copyright 2021 The Kruise Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package router

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"

	ctrlmeshutil "github.com/openkruise/kruise/pkg/ctrlmesh/util"
	"github.com/openkruise/kruise/pkg/util"
	"github.com/openkruise/kruise/pkg/util/bufferpool"
	"golang.org/x/net/http2"
	"golang.org/x/sync/semaphore"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured/unstructuredscheme"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

const (
	hugeBufferBytes = 4 * 1024 * 1024
	maxBufferBytes  = 16 * 1024 * 1024
)

var (
	hugeListSemaphore = semaphore.NewWeighted(2)
)

var (
	builtInScheme                    = runtime.NewScheme()
	typedNegotiatedSerializer        = serializer.NewCodecFactory(builtInScheme)
	unstructuredNegotiatedSerializer = unstructuredscheme.NewUnstructuredNegotiatedSerializer()

	typedSerializersMap        = map[string]*rest.Serializers{}
	unstructuredSerializersMap = map[string]*rest.Serializers{}
)

func init() {
	utilruntime.Must(metav1.AddMetaToScheme(builtInScheme))
	metav1.AddToGroupVersion(builtInScheme, schema.GroupVersion{Version: "v1"})
	utilruntime.Must(clientgoscheme.AddToScheme(builtInScheme))
	utilruntime.Must(admissionv1beta1.AddToScheme(builtInScheme))

	parseSerializers(typedNegotiatedSerializer, typedSerializersMap)
	parseSerializers(unstructuredNegotiatedSerializer, unstructuredSerializersMap)
}

func parseSerializers(negotiatedSerializer runtime.NegotiatedSerializer, serializersMap map[string]*rest.Serializers) {
	internalGV := schema.GroupVersions{
		// always include the legacy group as a decoding target to handle non-error `Status` return types
		{
			Group:   "",
			Version: runtime.APIVersionInternal,
		},
	}

	serializerInfos := negotiatedSerializer.SupportedMediaTypes()
	for i := range serializerInfos {
		info := serializerInfos[i]

		s := &rest.Serializers{
			Encoder: info.Serializer,
			Decoder: negotiatedSerializer.DecoderToVersion(info.Serializer, internalGV),

			RenegotiatedDecoder: func(contentType string, params map[string]string) (runtime.Decoder, error) {
				info, ok := runtime.SerializerInfoForMediaType(serializerInfos, contentType)
				if !ok {
					return nil, fmt.Errorf("serializer for %s not registered", contentType)
				}
				return negotiatedSerializer.DecoderToVersion(info.Serializer, internalGV), nil
			},
		}
		if info.StreamSerializer != nil {
			s.StreamingSerializer = info.StreamSerializer.Serializer
			s.Framer = info.StreamSerializer.Framer
		}

		serializersMap[info.MediaType] = s
	}
}

func getMediaType(resp *http.Response) (string, error) {
	contentType := resp.Header.Get("Content-Type")
	if len(contentType) == 0 {
		return "", fmt.Errorf("no mediaType in response %s", util.DumpJSON(resp.Header))
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", fmt.Errorf("invalid mediaType %s in response %s", mediaType, util.DumpJSON(resp.Header))
	}
	return mediaType, nil
}

func newResponseSerializer(resp *http.Response, apiResource *metav1.APIResource, isStreaming bool) *responseSerializer {
	mediaType, err := getMediaType(resp)
	if err != nil {
		return nil
	}

	var ser *rest.Serializers
	switch mediaType {
	case runtime.ContentTypeJSON:
		ser = unstructuredSerializersMap[runtime.ContentTypeJSON]
	case runtime.ContentTypeProtobuf:
		gvk := schema.GroupVersionKind{Group: apiResource.Group, Version: apiResource.Version, Kind: apiResource.Kind}
		if !builtInScheme.Recognizes(gvk) {
			klog.Errorf("Response with mediaType=%s but resource %v not found in built-in scheme", mediaType, gvk)
			return nil
		}
		ser = typedSerializersMap[runtime.ContentTypeProtobuf]
	case runtime.ContentTypeYAML:
		if isStreaming {
			klog.Errorf("Response with mediaType=%s can not support streaming watch", mediaType)
			return nil
		}
		ser = unstructuredSerializersMap[runtime.ContentTypeYAML]
	default:
		klog.Errorf("Response with unknown mediaType=%s", mediaType)
		return nil
	}

	rs := &responseSerializer{
		mediaType:   mediaType,
		apiResource: apiResource,
		framer:      ser.Framer,
	}

	if isStreaming {
		// do not give the bytes back to pool
		rs.buf = bytes.NewBuffer(bufferpool.BytesPool.Get())
		rs.buf2 = bytes.NewBuffer(bufferpool.BytesPool.Get())
	} else {
		if resp.ContentLength > hugeBufferBytes {
			rs.bufPool = bufferpool.BytesHugePool
		} else {
			rs.bufPool = bufferpool.BytesPool
		}
		rs.buf = bytes.NewBuffer(rs.bufPool.Get())
		rs.buf2 = bytes.NewBuffer(rs.bufPool.Get())
	}

	contentEncoding := resp.Header.Get("Content-Encoding")
	switch contentEncoding {
	case "gzip":
		rs.gzipReader = bufferpool.GzipReaderPool.Get().(*gzip.Reader)
		rs.gzipWriter = bufferpool.GzipWriterPool.Get().(*gzip.Writer)
	}

	bodyReader := resp.Body
	if rs.gzipReader != nil {
		rs.gzipReader.Reset(bodyReader)
		bodyReader = rs.gzipReader
	}
	if isStreaming {
		bodyReader = ser.Framer.NewFrameReader(bodyReader)
	}
	rs.reader = &fullReader{reader: bodyReader}

	return rs
}

type serializerReaderCloser struct {
	io.Reader
	serializer *responseSerializer
}

func (s *serializerReaderCloser) Close() error {
	s.serializer.Release()
	return nil
}

type responseSerializer struct {
	mediaType   string
	apiResource *metav1.APIResource
	reader      *fullReader
	framer      runtime.Framer

	bufPool    ctrlmeshutil.BufferPool
	buf        *bytes.Buffer
	buf2       *bytes.Buffer
	gzipReader *gzip.Reader
	gzipWriter *gzip.Writer
}

func (s *responseSerializer) Release() {
	if s.bufPool != nil {
		s.bufPool.Put(s.buf.Bytes()[:s.buf.Cap()])
		s.bufPool.Put(s.buf2.Bytes()[:s.buf2.Cap()])
	}
	if s.gzipReader != nil {
		s.gzipReader.Close()
		bufferpool.GzipReaderPool.Put(s.gzipReader)
	}
	if s.gzipWriter != nil {
		s.gzipWriter.Close()
		s.gzipWriter.Reset(nil)
		bufferpool.GzipWriterPool.Put(s.gzipWriter)
	}
}

func (s *responseSerializer) reset() {
	s.buf.Reset()
	s.buf2.Reset()
}

func (s *responseSerializer) EncodeList(obj *unstructured.Unstructured) (io.ReadCloser, int, error) {
	s.reset()
	bodyBuf := s.buf
	outputBuf := s.buf2

	var w io.Writer = outputBuf
	if s.gzipWriter != nil {
		s.gzipWriter.Reset(outputBuf)
		w = s.gzipWriter
	}

	if s.mediaType == runtime.ContentTypeProtobuf {
		if err := s.encodeUnstructuredObject(obj, bodyBuf); err != nil {
			return nil, 0, err
		}

		gvk := schema.GroupVersionKind{Group: s.apiResource.Group, Version: s.apiResource.Version, Kind: s.apiResource.Kind + "List"}
		realObj, _ := builtInScheme.New(gvk)
		if err := s.convertJSONToProtobuf(bodyBuf.Bytes(), realObj, w); err != nil {
			return nil, 0, err
		}
	} else {
		if err := s.encodeUnstructuredObject(obj, w); err != nil {
			return nil, 0, err
		}
	}

	if s.gzipWriter != nil {
		s.gzipWriter.Close()
	}
	return &serializerReaderCloser{Reader: outputBuf, serializer: s}, outputBuf.Len(), nil
	//return outputBuf.Bytes(), nil
}

func (s *responseSerializer) EncodeWatch(e *metav1.WatchEvent) ([]byte, error) {
	s.reset()
	outputBuf := s.buf

	if s.mediaType == runtime.ContentTypeProtobuf {
		if err := s.encodeUnstructuredObject(e.Object.Object.(*unstructured.Unstructured), outputBuf); err != nil {
			return nil, err
		}

		gvk := schema.GroupVersionKind{Group: s.apiResource.Group, Version: s.apiResource.Version, Kind: s.apiResource.Kind}
		realObj, _ := builtInScheme.New(gvk)
		if err := s.convertJSONToProtobuf(outputBuf.Bytes(), realObj, s.buf2); err != nil {
			return nil, err
		}
		e.Object.Object = nil
		e.Object.Raw = s.buf2.Bytes()
		outputBuf.Reset()
	}

	var w = s.framer.NewFrameWriter(outputBuf)
	if s.gzipWriter != nil {
		s.gzipWriter.Reset(w)
		w = s.gzipWriter
	}
	if err := s.encodeWatchEvent(e, w); err != nil {
		return nil, err
	}
	if s.gzipWriter != nil {
		s.gzipWriter.Close()
	}
	return outputBuf.Bytes(), nil
}

func (s *responseSerializer) DecodeList() (*unstructured.Unstructured, error) {
	s.reset()
	body, err := s.reader.readOnce()
	if err != nil {
		return nil, err
	}
	logBody("Response Body", body)

	if s.mediaType == runtime.ContentTypeProtobuf {
		gvk := schema.GroupVersionKind{Group: s.apiResource.Group, Version: s.apiResource.Version, Kind: s.apiResource.Kind + "List"}
		realObj, _ := builtInScheme.New(gvk)
		if err = s.convertProtobufToJSON(body, realObj, s.buf); err != nil {
			return nil, err
		}
		body = s.buf.Bytes()
	}

	return s.decodeUnstructuredObject(body)
}

func (s *responseSerializer) DecodeWatch() (*metav1.WatchEvent, error) {
	s.reset()
	n, err := s.reader.readStreaming(s.buf.Bytes()[:s.buf.Cap()])
	if err != nil {
		return nil, err
	}

	event, err := s.decodeWatchEvent(s.buf.Bytes()[:n])
	if err != nil {
		return nil, err
	}
	jsonBody := event.Object.Raw

	if s.mediaType == runtime.ContentTypeProtobuf {
		gvk := schema.GroupVersionKind{Group: s.apiResource.Group, Version: s.apiResource.Version, Kind: s.apiResource.Kind}
		realObj, _ := builtInScheme.New(gvk)
		if err = s.convertProtobufToJSON(event.Object.Raw, realObj, s.buf2); err != nil {
			return nil, err
		}
		jsonBody = s.buf2.Bytes()
	}

	obj, err := s.decodeUnstructuredObject(jsonBody)
	if err != nil {
		return nil, err
	}
	return &metav1.WatchEvent{Type: event.Type, Object: runtime.RawExtension{Object: obj}}, nil
}

func (s *responseSerializer) encodeWatchEvent(e *metav1.WatchEvent, w io.Writer) error {
	var encoder runtime.Encoder
	switch s.mediaType {
	case runtime.ContentTypeProtobuf:
		encoder = typedSerializersMap[runtime.ContentTypeProtobuf].StreamingSerializer
	default:
		encoder = typedSerializersMap[runtime.ContentTypeJSON].StreamingSerializer
	}
	return encoder.Encode(e, w)
}

func (s *responseSerializer) decodeWatchEvent(body []byte) (*metav1.WatchEvent, error) {
	e := &metav1.WatchEvent{}
	var decoder runtime.Decoder
	switch s.mediaType {
	case runtime.ContentTypeProtobuf:
		decoder = typedSerializersMap[runtime.ContentTypeProtobuf].StreamingSerializer
	default:
		decoder = typedSerializersMap[runtime.ContentTypeJSON].StreamingSerializer
	}
	if _, _, err := decoder.Decode(body, nil, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *responseSerializer) encodeUnstructuredObject(obj *unstructured.Unstructured, w io.Writer) error {
	var encoder runtime.Encoder
	switch s.mediaType {
	case runtime.ContentTypeProtobuf, runtime.ContentTypeJSON:
		encoder = unstructuredSerializersMap[runtime.ContentTypeJSON].Encoder
	default:
		encoder = unstructuredSerializersMap[runtime.ContentTypeYAML].Encoder
	}
	return encoder.Encode(obj, w)
}

func (s *responseSerializer) decodeUnstructuredObject(body []byte) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}
	var decoder runtime.Decoder
	switch s.mediaType {
	case runtime.ContentTypeProtobuf, runtime.ContentTypeJSON:
		decoder = unstructuredSerializersMap[runtime.ContentTypeJSON].Decoder
	default:
		decoder = unstructuredSerializersMap[runtime.ContentTypeYAML].Decoder
	}
	out, _, err := decoder.Decode(body, nil, obj)
	if err != nil {
		return nil, err
	} else if out == obj {
		return obj, nil
	}
	// if a different object is returned, see if it is Status and avoid double decoding
	// the object.
	switch t := out.(type) {
	case *metav1.Status:
		// any status besides StatusSuccess is considered an error.
		if t.Status != metav1.StatusSuccess {
			return nil, errors.FromObject(t)
		}
	}
	return obj, nil
}

func (s *responseSerializer) convertJSONToProtobuf(body []byte, obj runtime.Object, w io.Writer) error {
	jsonDecoder := typedSerializersMap[runtime.ContentTypeJSON].Decoder
	if _, _, err := jsonDecoder.Decode(body, nil, obj); err != nil {
		return err
	}

	protobufEncoder := typedSerializersMap[runtime.ContentTypeProtobuf].Encoder
	if err := protobufEncoder.Encode(obj, w); err != nil {
		klog.Errorf("Failed to encode object %s into protobuf", util.DumpJSON(obj))
		return fmt.Errorf("failed to encode into protobuf")
	}
	return nil
}

func (s *responseSerializer) convertProtobufToJSON(body []byte, obj runtime.Object, w io.Writer) error {
	protobufDecoder := typedSerializersMap[runtime.ContentTypeProtobuf].Decoder
	out, _, err := protobufDecoder.Decode(body, nil, obj)
	if err != nil {
		return err
	} else if out != obj {
		// if a different object is returned, see if it is Status and avoid double decoding
		// the object.
		switch t := out.(type) {
		case *metav1.Status:
			// any status besides StatusSuccess is considered an error.
			if t.Status != metav1.StatusSuccess {
				return errors.FromObject(t)
			}
		}
	}

	jsonEncoder := typedSerializersMap[runtime.ContentTypeJSON].Encoder
	if err := jsonEncoder.Encode(obj, w); err != nil {
		klog.Errorf("Failed to encode object %s into json", util.DumpJSON(obj))
		return fmt.Errorf("failed to encode into json")
	}
	return nil
}

// logBody logs a body output that could be either JSON or protobuf. It explicitly guards against
// allocating a new string for the body output unless necessary. Uses a simple heuristic to determine
// whether the body is printable.
func logBody(prefix string, body []byte) {
	if klog.V(8) {
		if bytes.IndexFunc(body, func(r rune) bool {
			return r < 0x0a
		}) != -1 {
			klog.Infof("%s:\n%s", prefix, truncateBody(hex.Dump(body)))
		} else {
			klog.Infof("%s: %s", prefix, truncateBody(string(body)))
		}
	}
}

// truncateBody decides if the body should be truncated, based on the glog Verbosity.
func truncateBody(body string) string {
	max := 0
	switch {
	case bool(klog.V(10)):
		return body
	case bool(klog.V(9)):
		max = 10240
	case bool(klog.V(8)):
		max = 1024
	}

	if len(body) <= max {
		return body
	}

	return body[:max] + fmt.Sprintf(" [truncated %d chars]", len(body)-max)
}

type fullReader struct {
	reader    io.ReadCloser
	resetRead bool
}

var ErrObjectTooLarge = fmt.Errorf("object to decode was longer than maximum allowed size")

func (f *fullReader) readOnce() ([]byte, error) {
	body, err := ioutil.ReadAll(f.reader)
	f.reader.Close()
	switch err.(type) {
	case nil:
	case http2.StreamError:
		// This is trying to catch the scenario that the server may close the connection when sending the
		// response body. This can be caused by server timeout due to a slow network connection.
		// TODO: Add test for this. Steps may be:
		// 1. client-go (or kubectl) sends a GET request.
		// 2. Apiserver sends back the headers and then part of the body
		// 3. Apiserver closes connection.
		// 4. client-go should catch this and return an error.
		klog.V(2).Infof("Stream error %#v when reading response body, may be caused by closed connection.", err)
		return nil, fmt.Errorf("stream error when reading response body, may be caused by closed connection. Please retry. Original error: %v", err)
	default:
		klog.Errorf("Unexpected error when reading response body: %v", err)
		return nil, fmt.Errorf("unexpected error when reading response body. Please retry. Original error: %v", err)
	}

	return body, nil
}

func (f *fullReader) readStreaming(buf []byte) (int, error) {
	base := 0
	for {
		n, err := f.reader.Read(buf[base:])
		if err == io.ErrShortBuffer {
			if n == 0 {
				return base, fmt.Errorf("got short buffer with n=0, base=%d, cap=%d", base, cap(buf))
			}
			if f.resetRead {
				continue
			}
			// double the buffer size up to maxBytes
			if len(buf) < maxBufferBytes {
				base += n
				buf = append(buf, make([]byte, len(buf))...)
				continue
			}
			// must read the rest of the frame (until we stop getting ErrShortBuffer)
			f.resetRead = true
			base = 0
			return base, ErrObjectTooLarge
		}
		if err != nil {
			return base, err
		}
		if f.resetRead {
			// now that we have drained the large read, continue
			f.resetRead = false
			continue
		}
		base += n
		break
	}
	return base, nil
}
