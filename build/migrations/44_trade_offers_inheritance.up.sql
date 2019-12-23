ALTER TABLE trade__resource_offers RENAME TO trade__offers;

ALTER TABLE trade__offers ALTER COLUMN price TYPE INT USING price::integer; 