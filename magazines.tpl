private _arsenal_magazine_types = ["{{ StringsJoin . "\", \"" }}"];
private _arsenal_magazines = createHashMap;

{{ range . }}
private _mags_{{ . }}_nocamo = [

];
private _mags_{{ . }}_black = [

];
private _mags_{{ . }}_desert = [

];
private _mags_{{ . }}_woodland = [

];

private _{{ . }}_mags = createHashMapFromArray  [["nocamo", _mags_{{ . }}_nocamo], ["black", _mags_{{ . }}_black], ["desert", _mags_{{ . }}_desert], ["woodland", _mags_{{ . }}_woodland]];
_arsenal_magazines set ["{{ . }}", _{{ . }}_mags];
{{ end }}

_arsenal_magazines set ["types", _arsenal_magazine_types];

_arsenal_magazines
