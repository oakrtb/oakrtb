package com.oakrtb.sdk.codec;

import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.Message;
import com.google.protobuf.util.JsonFormat;
import com.oakrtb.openrtb.v2.BidRequest;
import com.oakrtb.openrtb.v2.BidResponse;

import java.nio.charset.StandardCharsets;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.node.ObjectNode;
import com.google.protobuf.Descriptors;
import java.io.IOException;

/**
 * Conversion utilities for OpenRTB protobuf models and JSON.
 *
 * <p>Emits numeric enums and preserves proto field names for serialization/deserialization aligned with the OpenRTB JSON specification.
 */
public final class Json {
  private static final JsonFormat.Printer PRINTER =
      JsonFormat.printer()
          .omittingInsignificantWhitespace()
          .printingEnumsAsInts()
          .preservingProtoFieldNames();

  private static final JsonFormat.Parser PARSER =
      JsonFormat.parser().ignoringUnknownFields();

  private static final ObjectMapper MAPPER = new ObjectMapper()
      .enable(DeserializationFeature.USE_BIG_DECIMAL_FOR_FLOATS)
      .enable(DeserializationFeature.FAIL_ON_TRAILING_TOKENS);
  private Json() {}
  private static String transform(String raw, Descriptors.Descriptor descriptor, boolean encode) throws IOException {
    JsonNode root = MAPPER.readTree(raw);
    adapt(root, descriptor, encode);
    return MAPPER.writeValueAsString(root);
  }
  private static void adapt(JsonNode node, Descriptors.Descriptor descriptor, boolean encode) throws IOException {
    if (node == null || !node.isObject()) throw new IOException("codec: expected JSON object");
    ObjectNode object = (ObjectNode) node;
    for (var field : descriptor.getFields()) {
      String key = object.has(field.getName()) ? field.getName() : field.getJsonName();
      JsonNode value = object.get(key);
      if (value == null || value.isNull()) continue;
      if (field.getName().equals("ext")) {
        if (encode) value = MAPPER.readTree(value.textValue());
        if (value == null || !value.isObject()) throw new IOException("codec: ext must be an object");
        if (encode) object.set(key, value); else object.put(key, MAPPER.writeValueAsString(value));
      } else if (field.getJavaType() == Descriptors.FieldDescriptor.JavaType.MESSAGE) {
        if (field.isRepeated()) {
          if (!value.isArray()) throw new IOException("codec: expected array");
          for (var item : value) adapt(item, field.getMessageType(), encode);
        } else adapt(value, field.getMessageType(), encode);
      } else if (field.getJavaType() == Descriptors.FieldDescriptor.JavaType.STRING) {
        if (field.isRepeated()) {
          if (!value.isArray()) throw new IOException("codec: expected string array");
          for (var item : value) if (!item.isTextual()) throw new IOException("codec: expected string");
        } else if (!value.isTextual()) throw new IOException("codec: expected string");
      } else if (field.getJavaType() == Descriptors.FieldDescriptor.JavaType.DOUBLE || field.getJavaType() == Descriptors.FieldDescriptor.JavaType.INT || field.getJavaType() == Descriptors.FieldDescriptor.JavaType.ENUM) {
        if (field.isRepeated()) {
          if (!value.isArray()) throw new IOException("codec: expected array");
          for (var item : value) if (!item.isNumber() || !Double.isFinite(item.doubleValue())) throw new IOException("codec: expected finite number");
        } else if (!value.isNumber() || !Double.isFinite(value.doubleValue())) throw new IOException("codec: expected finite number");
      }
    }
  }
  private static String input(String raw, Descriptors.Descriptor descriptor) throws InvalidProtocolBufferException {
    try { return transform(raw, descriptor, false); }
    catch (IOException e) { throw new InvalidProtocolBufferException(e); }
  }

  /**
   * Serializes a protobuf message as an OpenRTB JSON string.
   *
   * @param message message to serialize
   * @return a JSON string
   * @throws IllegalArgumentException if serialization fails
   */
  public static String toJson(Message message) {
    try {
      return transform(PRINTER.print(message), message.getDescriptorForType(), true);
    } catch (IOException e) {
      throw new IllegalArgumentException("protobuf json encode failed", e);
    }
  }

  /**
   * Serializes a protobuf message as UTF-8 JSON bytes.
   *
   * @param message message to serialize
   * @return JSON bytes
   */
  public static byte[] toJsonBytes(Message message) {
    return toJson(message).getBytes(StandardCharsets.UTF_8);
  }

  /**
   * Parses a {@link BidRequest} from a JSON string.
   *
   * @param json OpenRTB BidRequest JSON
   * @return the parsed request
   * @throws InvalidProtocolBufferException if JSON cannot be merged into the proto
   */
  public static BidRequest parseBidRequest(String json) throws InvalidProtocolBufferException {
    BidRequest.Builder b = BidRequest.newBuilder();
    PARSER.merge(input(json, b.getDescriptorForType()), b);
    return b.build();
  }

  /**
   * Parses a {@link BidResponse} from a JSON string.
   *
   * @param json OpenRTB BidResponse JSON
   * @return the parsed response
   * @throws InvalidProtocolBufferException if JSON cannot be merged into the proto
   */
  public static BidResponse parseBidResponse(String json) throws InvalidProtocolBufferException {
    BidResponse.Builder b = BidResponse.newBuilder();
    PARSER.merge(input(json, b.getDescriptorForType()), b);
    return b.build();
  }

}
