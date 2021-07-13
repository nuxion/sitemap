User-agent: Googlebot
Disallow: /nogooglebot/

User-agent: *
{{range $val := .Routes -}}
{{ if $val.Robot -}}
Allow: /{{ $val.Loc }}
{{ else -}}
Disallow: /{{ $val.Loc }}
{{- end }}
{{- end }}
    
Sitemap: {{ .FullURL }}{{ .SitemapPath }}sitemap.xml