private [
  {{ range $i, $role := .Roles -}}
  "_{{ $role.Name }}",
  {{- end }}
  {{ range $n, $s := .Squads -}}
  "_{{ $n }}",
  {{- end }}
  "_squads_inf", "_squads_at", "_squads_aa", "_squads_spec_inf", "_base_soldier_class", "_roles", "_squads"
];

_base_soldier_class = "FIXME";
/*
 * Soldier roles definition
role = [
    name str,
    [
        pri_weapon_class str,
        # mags int,
        optics str["None", "holo", "combat", "sniper"],
        muzzle_device bool,
        frontgrip_device bool,
    ],
    grenades,
    vest_type str["None", "normal", "heavy"],
    backpack bool,
    trait str["tl", "medic", "gl", "lat", "hat", "aa", "mg_assistant", "trained", "specops"]
];

Traits:
    * tl - grenades are replaces with smoke, uses sidearm
    * medic - backpack is filled with medical supplies
    * gl - grenades are replaced with GL ammunition
    * lat - light AT as secondary weapon
    * hat - heavy AT as secondary weapon
    * aa - AA as secondary weapon
    * mg_assistant - backpack is filled with MG bearer ammunition
    * crew - doesn't wear vest. Lightly armed.
    * trained - skill is increased by 20%
    * specops - skill is increased by 50%
 */

{{ range $i, $role := .Roles }}
_{{ $role.Name }} = ["{{ $role.Name }}", ["{{ $role.WpnClass }}", {{ $role.NMags }}, "{{ $role.Optics }}", {{ $role.MuzzleDevice }}, {{ $role.RailDevice }}], {{ $role.NGrenades}}, "{{ $role.VestType}}", {{ $role.Backpack }}, [{{ $role.Traits.Traits }}]];
{{- end }}

{{ range $sq_name, $sq := .Squads }}
 _{{ $sq_name }} = [
{{- range $i, $r := $sq -}}
_{{ $r }}{{if IsLastItem $i (len $sq)}}{{else}},{{end}}
{{- end -}}
];
{{- end }}

_squads_inf = [_base_sq, 1.5, _light_sq, 2, _heavy_sq, 0.5, _sniper_ft, 0.1]; // weigted squads for correct randomization. Lower the weight - lower spawn chance
_squads_at = [_at_sq, 1];
_squads_aa = [_aa_sq, 1];
_squads_spec = [_spec_ft, 1, _spec_sniper_ft, 0.2];

_roles = createHashMapFromArray [
    {{- $trailing_item := false }}
    {{ range $i, $role := .Roles }}
    {{- if $trailing_item }}, {{end}}["{{ $role.Name }}", _{{ $role.Name }}]{{ $trailing_item = true}}
    {{- end }}
];

_squads = createHashMapFromArray [
            ["regular", [_squads_inf, _squads_at, _squads_aa]],    // conventional squads [anti-infantry, anti-tank, anti-air]
            ["patrol", _patrol_sq],
            ["guard", _sentry_comp],
            ["ambush", _ambush_sq],
            ["crew", _crew_comp]
];

[
    _base_soldier_class,     // base class for spawning and re-arming, needed for correct voice and visual appearance
    _roles,
    _squads
]