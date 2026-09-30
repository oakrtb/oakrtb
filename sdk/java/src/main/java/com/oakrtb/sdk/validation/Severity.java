package com.oakrtb.sdk.validation;

/**
 * Severity of a {@link CheckIssue}.
 */
public enum Severity {
  /** Structural mismatch: callers should reject or rewrite the bid. */
  ERROR,
  /** Business risk: the payload may still be valid; log or filter according to policy. */
  WARN
}
