package help

const HELP = `
Welcome to my MUD! Here, you can explore, battle, and interact with other players. Use commands to navigate through the world, manage your equipment, and engage in combat. Below are some of the basic commands you'll need to get started.

This game is in very early alpha stage. Characters are NOT currently persisted.

General Commands:
  - help: Display this help message.
  - look: Look around the current room to see its description and any items or characters present.
  - score: View your character's statistics and status.
  - gossip <message>: Send a message to all players in the game.

Movement Commands:
  - north, south, east, west, up, down: Move in the specified direction.

Inventory and Equipment Commands:
  - inventory: Display the items in your inventory.
  - equipment: Display the items you are currently wearing or wielding.
  - get <item>: Pick up an item from the current room.
  - drop <item>: Drop an item from your inventory into the current room.
  - wear <item>: Wear an item from your inventory.
  - wield <item>: Wield a weapon from your inventory.
  - remove <item>: Remove an item you are wearing or wielding.

Combat Commands:
  - kill <target>: Attack a specified target in the current room.

Utility Commands:
  - recall: Return to your recall point.
  - quit: Exit the game.

All commands may be abbreviated, e.g., l for look, e for east, etc.
`
