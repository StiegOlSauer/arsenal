private _vests_hash = createHashMap;

{{- range $v_type, $v_camo_map := . }}
{{- if ne $v_type "Starting"}}
private _vests_{{ $v_type | ToLower }} = createHashMap;

{{- range $v_camo, $v_array := $v_camo_map }}
private _vests_{{ $v_type | ToLower }}_{{ $v_camo }} = [
    {{- range $i, $v := $v_array }}
    ["{{ $v.Name }}", "{{ $v.Camoflage }}"],     {{ $v.Score }}{{if IsLastItem $i (len $v_array)}}{{else}},{{end}}
    {{- end -}}
];

_vests_{{ $v_type | ToLower }} set ["{{ $v_camo}}", _vests_{{ $v_type | ToLower }}_{{ $v_camo }}];
{{ end }}
_vests_hash set ["{{ $v_type | ToLower }}", _vests_{{ $v_type | ToLower }}];
{{- end }}
{{ end }}

_vests_hash

