
create table if not exists userstate (
    chat_id integer primary key,
    user_name text,
	scenario_name text,
	step_name integer
                                     );