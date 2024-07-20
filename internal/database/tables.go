package database

const schema = `
	CREATE TABLE IF NOT EXISTS public."user"
	(
		id uuid NOT NULL DEFAULT gen_random_uuid(),
		created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
		cookie text COLLATE pg_catalog."default",
		ip text COLLATE pg_catalog."default",
		CONSTRAINT user_pkey PRIMARY KEY (id)
	);

	CREATE TABLE IF NOT EXISTS public.ingredient
	(
		id uuid NOT NULL DEFAULT gen_random_uuid(),
		created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
		name text COLLATE pg_catalog."default" NOT NULL,
		standard_unit text COLLATE pg_catalog."default",
		ndb_number bigint,
		category text COLLATE pg_catalog."default",
		fdic_id bigint,
		CONSTRAINT ingredient_pkey PRIMARY KEY (id),
		CONSTRAINT con_unique_name UNIQUE (name)
	);

	CREATE TABLE IF NOT EXISTS public.technique
	(
		id uuid NOT NULL DEFAULT gen_random_uuid(),
		created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
		description text COLLATE pg_catalog."default" NOT NULL,
		CONSTRAINT technique_pkey PRIMARY KEY (id)
	);

	CREATE TABLE IF NOT EXISTS public.recipes
	(
		id uuid NOT NULL DEFAULT gen_random_uuid(),
		created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
		author uuid NOT NULL,
		name text COLLATE pg_catalog."default" NOT NULL,
		cuisine text COLLATE pg_catalog."default",
		yield smallint,
		yield_unit text COLLATE pg_catalog."default",
		prep_time interval,
		cooking_time interval,
		version bigint DEFAULT 0,
		CONSTRAINT recipes_pkey PRIMARY KEY (id),
		CONSTRAINT fk_recipe_user FOREIGN KEY (author)
			REFERENCES public."user" (id) MATCH SIMPLE
			ON UPDATE CASCADE
			ON DELETE CASCADE
			NOT VALID
	);

	CREATE INDEX IF NOT EXISTS fki_fk_recipe_user
		ON public.recipes USING btree
		(author ASC NULLS LAST)
		TABLESPACE pg_default;

	CREATE TABLE IF NOT EXISTS public.recipe_ingredient
	(
		id uuid NOT NULL DEFAULT gen_random_uuid(),
		created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
		recipe_id uuid,
		ingredient_id uuid,
		amount bigint,
		unit text COLLATE pg_catalog."default",
		CONSTRAINT recipe_ingredient_pkey PRIMARY KEY (id),
		CONSTRAINT fk_recipe_ingredient_ingredient FOREIGN KEY (ingredient_id)
			REFERENCES public.ingredient (id) MATCH FULL
			ON UPDATE CASCADE
			ON DELETE NO ACTION,
		CONSTRAINT fk_recipe_ingredient_recipe FOREIGN KEY (recipe_id)
			REFERENCES public.recipes (id) MATCH SIMPLE
			ON UPDATE CASCADE
			ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS public.step
	(
		id uuid NOT NULL DEFAULT gen_random_uuid(),
		created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
		step text COLLATE pg_catalog."default",
		recipe_id uuid NOT NULL,
		technique_id uuid,
		ingredient_id uuid,
		CONSTRAINT step_pkey PRIMARY KEY (id),
		CONSTRAINT fk_step_recipe FOREIGN KEY (recipe_id)
			REFERENCES public.recipes (id) MATCH FULL
			ON UPDATE NO ACTION
			ON DELETE NO ACTION,
		CONSTRAINT fk_step_recipe_ingredient FOREIGN KEY (ingredient_id)
			REFERENCES public.recipe_ingredient (id) MATCH FULL
			ON UPDATE NO ACTION
			ON DELETE NO ACTION,
		CONSTRAINT fk_step_technique FOREIGN KEY (technique_id)
			REFERENCES public.technique (id) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE NO ACTION
	);

	CREATE TABLE IF NOT EXISTS public.nutritional_value
	(
		id uuid NOT NULL DEFAULT gen_random_uuid(),
		created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
		ingredient_id uuid,
		recipe_id uuid,
		kcal numeric,
		kj numeric,
		fat numeric,
		saturated_fat numeric,
		carbohydrate numeric,
		sugar numeric,
		protein numeric,
		salt numeric,
		CONSTRAINT nutritional_value_pkey PRIMARY KEY (id),
		CONSTRAINT nutritional_value_unique_fk_ingredient UNIQUE (ingredient_id),
		CONSTRAINT nutritional_value_unique_fk_recipe UNIQUE (recipe_id),
		CONSTRAINT fk_nutritional_value_ingredient FOREIGN KEY (ingredient_id)
			REFERENCES public.ingredient (id) MATCH SIMPLE
			ON UPDATE CASCADE
			ON DELETE CASCADE
			NOT VALID,
		CONSTRAINT fk_nutritional_value_recipe FOREIGN KEY (recipe_id)
			REFERENCES public.recipes (id) MATCH SIMPLE
			ON UPDATE CASCADE
			ON DELETE CASCADE
			NOT VALID,
		CONSTRAINT nutritional_value_check CHECK (((ingredient_id IS NOT NULL)::integer + (recipe_id IS NOT NULL)::integer) = 1)
	);

	CREATE TABLE IF NOT EXISTS public.rating
	(
		id uuid NOT NULL DEFAULT gen_random_uuid(),
		created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
		ingredient_id uuid,
		recipe_id uuid,
		overall numeric,
		mon numeric,
		tue numeric,
		wed numeric,
		thu numeric,
		fri numeric,
		sat numeric,
		sun numeric,
		win numeric,
		spr numeric,
		sum numeric,
		aut numeric,
		thirtydegree numeric,
		twentiedegree numeric,
		tendegree numeric,
		zerodegree numeric,
		subzerodegree numeric,
		CONSTRAINT rating_pkey PRIMARY KEY (id),
		CONSTRAINT rating_unique_fk_ingredient UNIQUE (ingredient_id),
		CONSTRAINT rating_unique_fk_recipe UNIQUE (recipe_id),
		CONSTRAINT fk_rating_ingredient FOREIGN KEY (ingredient_id)
			REFERENCES public.ingredient (id) MATCH SIMPLE
			ON UPDATE CASCADE
			ON DELETE CASCADE,
		CONSTRAINT fk_rating_recipe FOREIGN KEY (recipe_id)
			REFERENCES public.recipes (id) MATCH SIMPLE
			ON UPDATE CASCADE
			ON DELETE CASCADE,
		CONSTRAINT rating_check CHECK (((ingredient_id IS NOT NULL)::integer + (recipe_id IS NOT NULL)::integer) = 1)
	);
	CREATE TABLE IF NOT EXISTS public."user"
	(
		id uuid NOT NULL DEFAULT gen_random_uuid(),
		created_at timestamp without time zone DEFAULT (now())::timestamp without time zone,
		cookie text COLLATE pg_catalog."default",
		ip text COLLATE pg_catalog."default",
		CONSTRAINT user_pkey PRIMARY KEY (id)
	);

`
