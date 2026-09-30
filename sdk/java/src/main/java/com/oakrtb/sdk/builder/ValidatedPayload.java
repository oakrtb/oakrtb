package com.oakrtb.sdk.builder;

import com.oakrtb.sdk.jsonschema.Report;

/**
 * Result of {@link #buildValidated()}: OpenRTB JSON bytes and Schema validation results.
 *
 * @param json serialized JSON bytes
 * @param result Schema validation report
 */
public record ValidatedPayload(byte[] json, Report result) {
  /**
   * Reports whether validation passed ({@link Report#isOk()} is true).
   *
   * @return {@code true} on success
   */
  public boolean ok() {
    return result != null && result.isOk();
  }
}
