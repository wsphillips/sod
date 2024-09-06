import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../core/components/other_inputs.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { Class, Faction, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import * as ProtectionPaladinInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecProtectionPaladin, {
	cssClass: 'protection-paladin-sim-ui',
	cssScheme: 'paladin',
    // List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
		`Judgement of the Crusader is currently not implemented; users can manually award themselves the relevant spellpower amount
		for a dps gain that will be slightly inflated given JotC does not benefit from source damage modifiers.`,
		`Be aware that not all item and weapon enchants are currently implemented in the sim, which make some notable Retribution
		weapons like Pendulum of Doom and The Jackhammer undervalued.`,
	],
	warnings: [
		(simUI: IndividualSimUI<Spec.SpecProtectionPaladin>) => {
			return {
				updateOn: simUI.player.changeEmitter,
				getContent: () => {
					if (simUI.player.getSpecOptions().primarySeal == 0) {
						return `Your previously selected seal is no longer available because of a talent or rune change.
							No seal will be cast with this configuration. Please select an available seal in the Settings>Player menu.`;
					} else {
						return '';
					}
				},
			};
		},
	],
	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatSpellHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatSpellPower,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDefense,
		Stat.StatBlock,
		Stat.StatBlockValue,
		Stat.StatDodge,
		Stat.StatParry,
		Stat.StatResilience,
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatMeleeSpeedMultiplier],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatSpellPower,
		Stat.StatSpellHit,
		Stat.StatDefense,
		Stat.StatBlock,
		Stat.StatBlockValue,
		Stat.StatDodge,
		Stat.StatParry,
		Stat.StatResilience,
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatArmor]: 0.07,
				[Stat.StatBonusArmor]: 0.06,
				[Stat.StatStamina]: 1.14,
				[Stat.StatStrength]: 1.0,
				[Stat.StatAgility]: 0.62,
				[Stat.StatAttackPower]: 0.26,
				[Stat.StatMeleeHit]: 0.79,
				[Stat.StatMeleeCrit]: 0.3,
				[Stat.StatSpellPower]: 0.13,
				[Stat.StatBlock]: 0.52,
				[Stat.StatBlockValue]: 0.28,
				[Stat.StatDodge]: 0.46,
				[Stat.StatParry]: 0.61,
				[Stat.StatDefense]: 0.54,
                [Stat.StatFireResistance]: 0.5,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 3.33,
                [PseudoStat.PseudoStatMeleeSpeedMultiplier]: 3.33,
			},
		),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
        other: Presets.OtherDefaults,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: Presets.DefaultIndividualBuffs,
		debuffs: Presets.DefaultDebuffs,
        race: Race.RaceHuman,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [ProtectionPaladinInputs.PrimarySealSelection, ProtectionPaladinInputs.AuraSelection],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.SpellScorchDebuff],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.TankAssignment,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.BurstWindow,
			OtherInputs.HpPercentForDefensives,
			OtherInputs.InspirationUptime,
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		rotations: [
			...Presets.APLPresets[Phase.Phase4],
		],
		// Preset talents that the user can quickly select.
		talents: [
			...Presets.TalentPresets[Phase.Phase4],
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			...Presets.GearPresets[Phase.Phase4],
		],
	},

	autoRotation: player => {
		return Presets.DefaultAPLs[player.getLevel()].rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecProtectionPaladin,
			tooltip: 'Protection Paladin',
			defaultName: 'Protection',
			iconUrl: getSpecIcon(Class.ClassPaladin, 1),

			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceUnknown,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.GearPresets[Phase.Phase4][0].gear,
				},
				[Faction.Horde]: {},
			},
		},
	],
});

export class ProtectionPaladinSimUI extends IndividualSimUI<Spec.SpecProtectionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecProtectionPaladin>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
