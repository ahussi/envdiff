// Package templater generates a .env.template (or .env.example) file from
// one or more parsed environment files.
//
// It collects the union of all keys found across the provided env maps and
// writes them out with empty values (or optional placeholder hints), making
// it easy to share a safe template with collaborators.
//
// Usage:
//
//	envs := map[string]map[string]string{
//		".env.production": {"DB_HOST": "prod.db", "DB_PASSWORD": "s3cr3t"},
//		".env.staging":    {"DB_HOST": "stage.db", "LOG_LEVEL": "debug"},
//	}
//
//	entries := templater.Generate(envs, templater.Options{Placeholders: true})
//	templater.Write(os.Stdout, entries, templater.Options{Placeholders: true})
//
// Output:
//
//	DB_HOST=<url>
//	DB_PASSWORD=<secret>
//	LOG_LEVEL=<value>
package templater
