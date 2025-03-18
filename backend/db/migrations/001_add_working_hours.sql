-- Добавление поля working_hours в таблицу user
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS working_hours JSONB DEFAULT '{}';

-- Обновление существующих записей
UPDATE "user" SET working_hours = '{
    "monday": {"start": null, "end": null},
    "tuesday": {"start": null, "end": null},
    "wednesday": {"start": null, "end": null},
    "thursday": {"start": null, "end": null},
    "friday": {"start": null, "end": null},
    "saturday": {"start": null, "end": null},
    "sunday": {"start": null, "end": null}
}'::jsonb WHERE working_hours IS NULL; 