// Package cascader provides multi-layer .env resolution, simulating the
// cascade behaviour used by tools such as docker-compose and dotenv-flow.
//
// Layers are processed in order from lowest to highest priority. A key
// defined in a later layer overwrites the same key from an earlier layer.
// The Result records both the final resolved value and which layer was
// the winning source (provenance).
//
// Typical usage:
//
//	res, err := cascader.Cascade([]cascader.Layer{
//		{Label: ".env",       Path: ".env"},
//		{Label: ".env.prod",  Path: ".env.prod"},
//		{Label: ".env.local", Path: ".env.local"},
//	})
package cascader
