package ui

import (
	"codeberg.org/anaseto/gruid"
)

// Combat-related colors
var (
	// Attack colors
	ColorPlayerAttack  gruid.Color // Player attacks monsters
	ColorEnemyAttack   gruid.Color // Monsters attack player
	ColorNeutralAttack gruid.Color // Monster attacks another monster

	// Hit flash effect
	ColorHitFlash gruid.Color // Visual feedback when entity is hit

	// Status colors
	ColorDeath  gruid.Color // For death messages
	ColorCorpse gruid.Color // For corpse messages
)

func init() {
	// Initialize combat colors
	ColorPlayerAttack = ColorBlue    // Same as player color
	ColorEnemyAttack = ColorRed      // Same as monster color
	ColorNeutralAttack = ColorYellow // Neutral color for monster-monster

	ColorHitFlash = ColorForegroundEmph    // Bright white flash when hit
	ColorDeath = ColorRed                  // Death messages are red
	ColorCorpse = ColorForegroundSecondary // Corpse messages are corpse color
}
