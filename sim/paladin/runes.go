package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (paladin *Paladin) ApplyRunes() {
	paladin.registerTheArtOfWar()
	paladin.registerSheathOfLight()
	paladin.registerGuardedByTheLight()
	paladin.registerShockAndAwe()
	paladin.registerRV()
	paladin.registerAegis()
	// "RuneHeadFanaticism" is handled in Exorcism, Holy Shock, SoC, and SoR
	paladin.registerFanaticism()

	// "RuneHeadWrath" is handled in Exorcism, Holy Shock, Consecration (and Holy Wrath once implemented)

	paladin.registerHammerOfTheRighteous()
	paladin.registerMalleableProtection()
	// "RuneWristImprovedHammerOfWrath" is handled Hammer of Wrath
	paladin.applyPurifyingPower()
}

func (paladin *Paladin) registerMalleableProtection() {

	if !paladin.hasRune(proto.PaladinRune_RuneWaistMalleableProtection) {
		return
	}

	defendersResolveAPAura := core.DefendersResolveAttackPower(paladin.GetCharacter())

	paladin.RegisterAura(core.Aura{
		Label:    "Malleable Protection",
		ActionID: core.ActionID{SpellID: 458318},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellCode != SpellCode_PaladinHolyShield {
				return
			}
			if stacks := int32(paladin.GetStat(stats.Defense)); stacks > 0 {
				defendersResolveAPAura.Activate(sim)
				if defendersResolveAPAura.GetStacks() != stacks {
					defendersResolveAPAura.SetStacks(sim, stacks)
				}
			}
		},
	})
}

func (paladin *Paladin) registerAegis() {

	if !paladin.hasRune(proto.PaladinRune_RuneChestAegis) {
		return
	}

	// Aegis bonus applies to equipped SBV only and not base SBV gained via
	// strength. The bonus is additive with Shield Specialization.
	paladin.AddStat(stats.BlockValue, paladin.sbvEquipBonus(1.3))

	actionID := core.ActionID{SpellID: 425589}

	paladin.RegisterAura(core.Aura{
		Label:    "Aegis Rune",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				if sim.RandomFloat("Redoubt") < 0.1 {
					paladin.redoubtAura.Activate(sim)
					paladin.redoubtAura.SetStacks(sim, 5)
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMelee) && result.Outcome.Matches(core.OutcomeCrit) {
				paladin.redoubtAura.Activate(sim)
				paladin.redoubtAura.SetStacks(sim, 5)
			}
		},
	})

	// Reckoning now procs on any melee or ranged hit with (2% * talent points) chance
	reck_actionID := core.ActionID{SpellID: 20178} // reckoning proc id

	procChance := 0.02 * float64(paladin.Talents.Reckoning)

	paladin.RegisterAura(core.Aura{
		Label:    "Reckoning Aegis Proc Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) && sim.Proc(procChance, "Reckoning") {
				paladin.AutoAttacks.ExtraMHAttack(sim, 1, reck_actionID, spell.ActionID)
			}
		},
	})
}

func (paladin *Paladin) fanaticism() float64 {
	return core.TernaryFloat64(paladin.hasRune(proto.PaladinRune_RuneHeadFanaticism), 18, 0) * core.SpellCritRatingPerCritChance
}

func (paladin *Paladin) registerFanaticism() {
	if paladin.hasRune(proto.PaladinRune_RuneHeadFanaticism) {
		paladin.PseudoStats.SchoolBonusCritChance[stats.SchoolIndexHoly] += 18
	}
}

func (paladin *Paladin) registerTheArtOfWar() {
	if !paladin.hasRune(proto.PaladinRune_RuneFeetTheArtOfWar) {
		return
	}

	paladin.RegisterAura(core.Aura{
		Label:    "The Art of War",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee|core.ProcMaskTriggerInstant) || !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			//paladin.holyShockCooldown.Reset()
			paladin.exorcismCooldown.Set(sim.CurrentTime + max(0, paladin.exorcismCooldown.TimeToReady(sim)-(time.Second*2)))
		},
	})
}

func (paladin *Paladin) registerSheathOfLight() {

	if !paladin.hasRune(proto.PaladinRune_RuneWaistSheathOfLight) {
		return
	}

	dep := paladin.NewDynamicStatDependency(
		stats.AttackPower, stats.SpellPower, 0.3,
	)

	sheathAura := paladin.RegisterAura(core.Aura{
		Label:    "Sheath of Light",
		Duration: time.Second * 60,
		ActionID: core.ActionID{SpellID: 426159},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.EnableDynamicStatDep(sim, dep)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.DisableDynamicStatDep(sim, dep)
		},
	})
	paladin.RegisterAura(core.Aura{
		Label:    "Sheath of Light (rune)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskWhiteHit) {
				return
			}
			sheathAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) registerShockAndAwe() {

	if !paladin.hasRune(proto.PaladinRune_RuneCloakShockAndAwe) {
		return
	}

	dep := paladin.NewDynamicStatDependency(
		stats.Intellect, stats.SpellDamage, 1.0,
	)

	shockAndAweAura := paladin.RegisterAura(core.Aura{
		Label:    "Shock and Awe",
		Duration: time.Second * 60,
		ActionID: core.ActionID{SpellID: 462834},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.EnableDynamicStatDep(sim, dep)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.DisableDynamicStatDep(sim, dep)
		},
	})
	paladin.RegisterAura(core.Aura{
		Label:    "Shock and Awe (rune)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellCode != SpellCode_PaladinHolyShock {
				return
			}
			shockAndAweAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) registerGuardedByTheLight() {
	if !paladin.hasRune(proto.PaladinRune_RuneFeetGuardedByTheLight) {
		return
	}

	actionID := core.ActionID{SpellID: 415058}
	manaMetrics := paladin.NewManaMetrics(actionID)
	var manaPA *core.PendingAction

	guardedAura := paladin.RegisterAura(core.Aura{
		Label:    "Guarded by the Light",
		Duration: time.Second*15 + 1,
		ActionID: actionID,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			manaPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 3,
				OnAction: func(sim *core.Simulation) {
					paladin.AddMana(sim, 0.05*paladin.MaxMana(), manaMetrics)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			manaPA.Cancel(sim)
		},
	})

	paladin.RegisterAura(core.Aura{
		Label:    "Guarded by the Light (rune)",
		Duration: core.NeverExpires,
		ActionID: core.ActionID{SpellID: 415755},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskWhiteHit) {
				return
			}
			guardedAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) applyPurifyingPower() {
	if !paladin.hasRune(proto.PaladinRune_RuneWristPurifyingPower) {
		return
	}

	paladin.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellCode == SpellCode_PaladinExorcism || spell.SpellCode == SpellCode_PaladinHolyWrath {
			spell.CD.Duration /= 2
		}
	})
}
