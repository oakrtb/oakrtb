//! Serialization support for generated model fields; no SDK layer dependencies.
use serde::{Serialize, Serializer};
pub(crate) fn is_zero(v: &i32) -> bool {
    *v == 0
}
pub(crate) fn finite_optional<S: Serializer>(
    value: &Option<f64>,
    serializer: S,
) -> Result<S::Ok, S::Error> {
    if value.is_some_and(|v| !v.is_finite()) {
        return Err(serde::ser::Error::custom("non-finite number"));
    }
    match value {
        Some(v) if v.fract() == 0.0 && v.abs() <= 9_007_199_254_740_992.0 => {
            serializer.serialize_i64(*v as i64)
        }
        _ => value.serialize(serializer),
    }
}
// Preserve the existing protobuf wire string while exposing an object in OpenRTB JSON.
pub(crate) mod ext {
    use serde::{Deserialize, Deserializer, Serialize, Serializer};
    use serde_json::Value;
    pub fn serialize<S: Serializer>(raw: &str, serializer: S) -> Result<S::Ok, S::Error> {
        let value: Value = serde_json::from_str(raw).map_err(serde::ser::Error::custom)?;
        if !value.is_object() {
            return Err(serde::ser::Error::custom("ext must be an object"));
        }
        value.serialize(serializer)
    }
    pub fn deserialize<'de, D: Deserializer<'de>>(deserializer: D) -> Result<String, D::Error> {
        let value = Option::<Value>::deserialize(deserializer)?;
        match value {
            None => Ok(String::new()),
            Some(v) if v.is_object() => Ok(v.to_string()),
            _ => Err(serde::de::Error::custom("ext must be an object")),
        }
    }
}

pub(crate) fn default_null<'de, D, T>(de: D) -> Result<T, D::Error>
where
    D: serde::Deserializer<'de>,
    T: serde::Deserialize<'de> + Default,
{
    use serde::Deserialize;
    Ok(Option::<T>::deserialize(de)?.unwrap_or_default())
}
fn number_i32<E: serde::de::Error>(n: serde_json::Number) -> Result<i32, E> {
    // Work on decimal digits: conversion through f64 can round a fraction into an integer.
    let raw = n.to_string();
    let (mantissa, exponent) = raw.split_once(['e', 'E']).unwrap_or((&raw, "0"));
    let exponent = exponent.parse::<i64>().unwrap_or_else(|_| {
        if exponent.starts_with('-') {
            i64::MIN
        } else {
            i64::MAX
        }
    });
    let negative = mantissa.starts_with('-');
    let unsigned = mantissa.trim_start_matches('-');
    let (whole, fraction) = unsigned.split_once('.').unwrap_or((unsigned, ""));
    let mut digits = format!("{whole}{fraction}")
        .trim_start_matches('0')
        .to_owned();
    if digits.is_empty() {
        return Ok(0);
    }
    let scale = exponent.saturating_sub(fraction.len() as i64);
    if scale < 0 {
        let trim = scale.unsigned_abs();
        if trim >= digits.len() as u64 {
            return Err(E::custom("expected integer"));
        }
        let keep = digits.len() - trim as usize;
        if !digits[keep..].bytes().all(|b| b == b'0') {
            return Err(E::custom("expected integer"));
        }
        digits.truncate(keep);
    } else {
        if scale > 10 || digits.len() + scale as usize > 10 {
            return Err(E::custom("integer outside int32 range"));
        }
        digits.extend(std::iter::repeat('0').take(scale as usize));
    }
    if negative {
        digits.insert(0, '-');
    }
    digits
        .parse::<i32>()
        .map_err(|_| E::custom("integer outside int32 range"))
}
pub(crate) fn optional_i32<'de, D: serde::Deserializer<'de>>(
    de: D,
) -> Result<Option<i32>, D::Error> {
    use serde::Deserialize;
    Option::<serde_json::Number>::deserialize(de)?
        .map(number_i32)
        .transpose()
}
pub(crate) fn scalar_i32<'de, D: serde::Deserializer<'de>>(de: D) -> Result<i32, D::Error> {
    Ok(optional_i32(de)?.unwrap_or(0))
}
pub(crate) fn repeated_i32<'de, D: serde::Deserializer<'de>>(de: D) -> Result<Vec<i32>, D::Error> {
    let numbers: Vec<serde_json::Number> = default_null(de)?;
    numbers.into_iter().map(number_i32).collect()
}

pub(crate) fn optional_message<'de, D, T>(de: D) -> Result<Option<T>, D::Error>
where
    D: serde::Deserializer<'de>,
    T: serde::de::DeserializeOwned,
{
    use serde::Deserialize;
    let value = Option::<serde_json::Value>::deserialize(de)?;
    value
        .map(|v| {
            if !v.is_object() {
                return Err(serde::de::Error::custom("expected JSON object"));
            }
            serde_json::from_value(v).map_err(serde::de::Error::custom)
        })
        .transpose()
}
pub(crate) fn repeated_message<'de, D, T>(de: D) -> Result<Vec<T>, D::Error>
where
    D: serde::Deserializer<'de>,
    T: serde::de::DeserializeOwned,
{
    let values: Vec<serde_json::Value> = default_null(de)?;
    values
        .into_iter()
        .map(|v| {
            if !v.is_object() {
                return Err(serde::de::Error::custom("expected JSON object"));
            }
            serde_json::from_value(v).map_err(serde::de::Error::custom)
        })
        .collect()
}
