private _items_hash = createHashMap;

{{- range $i_type, $i_array := . }}
private _{{$i_type}} = [
    {{- $printed_once := false}}
    {{- range $i, $u := $i_array }}
    {{- if $printed_once }}, {{end}}"{{$u}}"{{ $printed_once = true}}
    {{- end -}}
];

_items_hash set ["{{ $i_type | ToLower }}", _{{ $i_type }}];
{{ end }}

_items_hash