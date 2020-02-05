-- truncate table bookms_user.bookms_user_comments;
SELECT * FROM bookms_user.bookms_user_comments;
-- insert into bookms_user.bookms_user_score(identify,user_id,score,create_time) values(2,3,50,now());
use bookms_user;
-- insert into bookms_user_comments(identify,user_id,content,create_time) values(2,3,"test",Now());
-- drop procedure if exists BatchInsert;
-- delimiter $ -- 把界定符改成双斜杠
-- create procedure BatchInsert(IN loop_count INT) 
-- begin
-- 	declare Var INT;
--     set Var = 0;
--     while Var < loop_count do
-- 		insert into bookms_user_comments(identify,user_id,content,create_time) values(2,3,'test',Now());
--         set Var = Var + 1;
-- 	end while;
-- end;
-- call BatchInsert(10);
-- delimiter ;
