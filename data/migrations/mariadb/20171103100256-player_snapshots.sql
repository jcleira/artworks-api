-- +migrate Up
CREATE TABLE artworks (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  rei VARCHAR(9) NOT NULL,
  created_at INT NOT NULL,
  ubi VARCHAR(255) NOT NULL,
  pro VARCHAR(255) NOT NULL,
  adq VARCHAR(255) NOT NULL,
  reg VARCHAR(255) NOT NULL,
  nom VARCHAR(255) NOT NULL,
  tit VARCHAR(255) NOT NULL,
  aut VARCHAR(255) NOT NULL,
  fec VARCHAR(255) NOT NULL,
  lug VARCHAR(255) NOT NULL,
  ico VARCHAR(255) NOT NULL,
  icc VARCHAR(255) NOT NULL,
  tip VARCHAR(255) NOT NULL,
  tec VARCHAR(255) NOT NULL,
  sop VARCHAR(255) NOT NULL,
  mat VARCHAR(255) NOT NULL,
  tin VARCHAR(255) NOT NULL,
  dim VARCHAR(255) NOT NULL,
  hue VARCHAR(255) NOT NULL,
  ins VARCHAR(255) NOT NULL,
  des VARCHAR(255) NOT NULL,
  est VARCHAR(255) NOT NULL,
  INDEX `id` (`id`)
);

-- +migrate Down
DROP TABLE artworks;