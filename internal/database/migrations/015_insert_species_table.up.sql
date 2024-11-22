CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Dados para Milho (Zea mays)
INSERT INTO species (
    id, common_name, specie_description, scientific_name, botanical_family, 
    growth_type, ideal_temperature, ideal_climate, life_cycle, planting_season, 
    harvest_time, average_height, average_width, irrigation_weight, 
    fertilization_weight, sun_weight, created_at, updated_at
) VALUES (
    gen_random_uuid(), -- ou substitua por um UUID gerado
    'Milho', 
    'Planta anual cultivada para grãos e forragem.', 
    'Zea mays', 
    'Poaceae', 
    'Herbácea', 
    21.0, 
    'Tropical e subtropical', 
    'Anual', 
    'Primavera/Verão', 
    120, 
    2.50, 
    0.25, 
    0.75, 
    0.60, 
    0.85, 
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP
);

-- Dados para Laranja (Citrus sinensis)
INSERT INTO species (
    id, common_name, specie_description, scientific_name, botanical_family, 
    growth_type, ideal_temperature, ideal_climate, life_cycle, planting_season, 
    harvest_time, average_height, average_width, irrigation_weight, 
    fertilization_weight, sun_weight, created_at, updated_at
) VALUES (
    gen_random_uuid(), -- ou substitua por um UUID gerado
    'Laranja', 
    'Árvore perene cultivada para frutos cítricos.', 
    'Citrus sinensis', 
    'Rutaceae', 
    'Árvores frutíferas', 
    24.0, 
    'Tropical e subtropical', 
    'Perene', 
    'Outono/Inverno', 
    365, 
    4.50, 
    2.50, 
    0.70, 
    0.65, 
    0.90, 
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP
);
