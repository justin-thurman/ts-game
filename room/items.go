package room

import (
	"fmt"
	"slices"
	"ts-game/items"
	"ts-game/player"
)

// HandleGet handles getting an item by name.
func (r *Room) HandleGet(p *player.Player, itemName string) (message string) {
	r.itemsMu.Lock()
	defer r.itemsMu.Unlock()
	var targetItem items.RoomItem
	for i, it := range r.roomItems {
		if it.HasName(itemName) {
			targetItem = it
			r.roomItems = slices.Delete(r.roomItems, i, i+1)
			break
		}
	}
	if targetItem == nil {
		return fmt.Sprintf("There doesn't seem to be a %s here.", itemName)
	}
	defer r.updateDescription()
	targetItem.Get(p.Inventory)
	return fmt.Sprintf("You pick up %s.", targetItem.String())
}

// HandleDrop handles dropping an item by name.
func (r *Room) HandleDrop(p *player.Player, itemName string) (message string) {
	item := p.Inventory.GetItemForDropping(itemName)
	if item == nil {
		return fmt.Sprintf("You don't seem to be carrying a %s.", itemName)
	}
	r.itemsMu.Lock()
	defer r.itemsMu.Unlock()
	defer r.updateDescription()
	r.roomItems = append(r.roomItems, item)
	return fmt.Sprintf("You drop %s.", item.String())
}
