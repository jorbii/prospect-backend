ALTER TABLE inventory_items
ADD CONSTRAINT inventory_items_item_id_fkey
FOREIGN KEY (item_id)
REFERENCES items(id)
ON DELETE RESTRICT;