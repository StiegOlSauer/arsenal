private _devs_hash = createHashMap;

{{- range $d_camo, $d_type_map := . }}
private _devs_{{ $d_camo }} = createHashMap;

{{- range $d_type, $d_array := $d_type_map }}
private _devs_{{ $d_camo }}_{{ $d_type | ToLower }} = [
    {{- range $i, $d := $d_array }}
    "{{ $d.Name }}",     {{ $d.Score }}{{if IsLastItem $i (len $d_array)}}{{else}},{{end}}
    {{- end -}}
];

_devs_{{ $d_camo }} set ["{{ $d_type | ToLower }}", _devs_{{ $d_camo }}_{{ $d_type | ToLower }}];
{{ end }}
_devs_hash set ["{{ $d_camo }}", _devs_{{ $d_camo }}];
{{- end }}

_devs_hash

