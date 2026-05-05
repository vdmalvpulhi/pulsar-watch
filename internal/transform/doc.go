// Package transform provides payload transformation support for pulsar-watch.
//
// A Transformer is created from a Config that specifies:
//
//   - MaskFields: JSON field names to redact (values replaced with "***").
//   - TruncateAt: maximum byte length before the payload is clipped.
//   - RenameFields: a map of old→new JSON field name substitutions.
//
// Transformations are applied in order: rename → mask → truncate.
// Non-JSON payloads pass through rename/mask unchanged and are only
// subject to truncation.
//
// Example:
//
//	cfg := transform.Config{
//		MaskFields:   []string{"password", "token"},
//		TruncateAt:   256,
//		RenameFields: map[string]string{"userId": "user_id"},
//	}
//	tr, err := transform.New(cfg)
//	if err != nil { ... }
//	result := tr.Apply(rawPayload)
package transform
