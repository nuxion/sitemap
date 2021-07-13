User-agent: Googlebot
Disallow: /nogooglebot/

User-agent: *
{{range $val := .Rules }}
{{ $val }}
{{- end }}
    
Sitemap: {{ .Sitemap }}