CREATE SEQUENCE IF NOT EXISTS rides_id_seq;

CREATE TABLE IF NOT EXISTS rides (
  id TEXT PRIMARY KEY DEFAULT ('ride-' || nextval('rides_id_seq')::TEXT),
  rider_id TEXT NOT NULL,
  pickup TEXT NOT NULL,
  dropoff TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
