package com.oakrtb.sdk.validation;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Objects;

/**
 * Aggregated bid check results. {@link #ok()} is true when no {@link Severity#ERROR} findings exist.
 */
public final class CheckResult {
  private final List<CheckIssue> issues;

  private CheckResult(List<CheckIssue> issues) {
    this.issues = List.copyOf(Objects.requireNonNull(issues, "issues"));
  }

  /**
   * Constructs results from a list of issues.
   *
   * @param issues findings
   * @return CheckResult
   */
  public static CheckResult of(List<CheckIssue> issues) {
    return new CheckResult(issues);
  }

  /**
   * Empty result with no issues.
   *
   * @return an empty CheckResult
   */
  public static CheckResult empty() {
    return new CheckResult(List.of());
  }

  /**
   * All findings.
   *
   * @return an immutable issue list
   */
  public List<CheckIssue> issues() {
    return issues;
  }

  /**
   * True when there are no ERROR findings (WARN findings may still exist).
   *
   * @return true on success
   */
  public boolean ok() {
    for (CheckIssue i : issues) {
      if (i.severity() == Severity.ERROR) {
        return false;
      }
    }
    return true;
  }

  /**
   * ERROR findings only.
   *
   * @return the ERROR list
   */
  public List<CheckIssue> errors() {
    List<CheckIssue> out = new ArrayList<>();
    for (CheckIssue i : issues) {
      if (i.severity() == Severity.ERROR) {
        out.add(i);
      }
    }
    return Collections.unmodifiableList(out);
  }

  /**
   * WARN findings only.
   *
   * @return the WARN list
   */
  public List<CheckIssue> warnings() {
    List<CheckIssue> out = new ArrayList<>();
    for (CheckIssue i : issues) {
      if (i.severity() == Severity.WARN) {
        out.add(i);
      }
    }
    return Collections.unmodifiableList(out);
  }

  /**
   * Reports whether the specified error code is present.
   *
   * @param code an {@link IssueCode} constant
   * @return true if present
   */
  public boolean has(String code) {
    for (CheckIssue i : issues) {
      if (code.equals(i.code())) {
        return true;
      }
    }
    return false;
  }
}
