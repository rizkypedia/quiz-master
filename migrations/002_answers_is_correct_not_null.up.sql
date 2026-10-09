UPDATE public.answers
SET is_correct = false
WHERE is_correct IS NULL;

ALTER TABLE public.answers
    ALTER COLUMN is_correct SET DEFAULT false,
    ALTER COLUMN is_correct SET NOT NULL;
