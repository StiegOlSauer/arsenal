# arsenal - Complementary faction configurator for Liberation mission for Arma 3

Generates SQF code from provided CSV tables which is then consumed by Arma 3 Liberation mission (https://github.com/StiegOlSauer/pclf_liberation.Altis). That allows to design whole faction by putting its properties into a spreadsheet instead of coding it in mission files. Liberation factions are quit complex, so the main benefits from this approach are: 
1. Integrity: it is only one file to edit instead of multitude of SQFs
2. Convenience: faction features are logically divided on several independent groups, each requires to fill several properties
3. Maintainability: there is no need to adapt faction file to changes in mission logic

Currently available entities:
* weapons
* optics
* body armor and headgear
* uniforms
* muzzle devices, pointers and frontgrip attachments
* explosives, mines, grenades
* vehicles
* squad roles and compositions

Usage example:
1. Prepare new or get existing .ods equipment file from Liberation repository (i.e. https://github.com/StiegOlSauer/pclf_liberation.Altis/blob/main/utils/arsenal/AFRFPMC.ods)
2. Export .csv files from its sheets into some convenient place (i.e. `arsenal/afrfpmc/weapons.csv`)
3. Run the generator from that directory:

  `$ ../arsenal -w --weapons weapons.csv --squads squads.csv --templates-dir /path/to/templates/`

4. SQF files in corresponding subdirs will be created in the directory from which `arsenal` is launched
5. Copy the directory into Liberation mission `scripts/gameplay_templates/arsenal/`
