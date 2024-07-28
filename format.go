package censor

// Format processes a supported value, converting it into a formatted string while adhering to secure data masking principles.
// This pivotal function within the censor package seamlessly handles various types, including structs, slices, arrays, pointers,
// maps, and etc., parsing them recursively.
//
// By default, struct fields undergo masking to safeguard sensitive data. For selective exposure, developers can utilize the
// `censor:"display"` tag, providing granular control over which fields remain visible.
//
// For bug reports or feedback, please contribute to the censor project at https://github.com/vpakhuchyi/censor.
func Format(val any) string {
	return globalInstance.Format(val)
}
