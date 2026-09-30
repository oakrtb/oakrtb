package com.oakrtb.sdk.jsonschema;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Objects;

/**
 * A single Schema validation failure.
 */
public final class Issue {
  private final String code;
  private final String path;
  private final String message;

  /**
   * @param code error category (e.g. required, type, format)
   * @param path JSON path
   * @param message detailed message
   */
  @JsonCreator
  public Issue(
      @JsonProperty("code") String code,
      @JsonProperty("path") String path,
      @JsonProperty("message") String message) {
    this.code = Objects.requireNonNull(code, "code");
    this.path = path == null ? "" : path;
    this.message = Objects.requireNonNull(message, "message");
  }

  @JsonProperty("code")
  public String getCode() {
    return code;
  }

  @JsonProperty("path")
  public String getPath() {
    return path;
  }

  @JsonProperty("message")
  public String getMessage() {
    return message;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof Issue)) {
      return false;
    }
    Issue that = (Issue) o;
    return Objects.equals(code, that.code)
        && Objects.equals(path, that.path)
        && Objects.equals(message, that.message);
  }

  @Override
  public int hashCode() {
    return Objects.hash(code, path, message);
  }
}
