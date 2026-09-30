package com.oakrtb.sdk.validation;

import java.util.Objects;

/**
 * A finding produced by validation.
 *
 * @param code stable error code (see {@link IssueCode})
 * @param severity severity level
 * @param path JSON-style path
 * @param message human-readable explanation
 */
public record CheckIssue(String code, Severity severity, String path, String message) {
  public CheckIssue {
    Objects.requireNonNull(code, "code");
    Objects.requireNonNull(severity, "severity");
    path = path == null ? "" : path;
    message = message == null ? "" : message;
  }
}
