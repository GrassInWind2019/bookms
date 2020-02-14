use bookms;
drop procedure if exists BatchInsert;
delimiter $ -- 把界定符改成双斜杠
create procedure BatchInsert(IN loop_count INT) 
begin
	declare Var INT;
    declare i INT;
    declare j int;
    declare k int;
    declare cid int;
    declare user_id int;
    declare lend_status int;
    set Var = 10001;
    set loop_count=loop_count+Var;
    -- select concat('loop_count=',loop_count);
    -- select concat('Var=',Var);
    loop1: while Var < loop_count do
		set lend_status = 0;
        set i = FLOOR(RAND()*1000+1);  -- 1-1000
		if i < 991 then
			set user_id = 1;
		elseif i < 1000 then
			set user_id = 3;
		else
			set user_id = 5;
		end if;
        if mod(Var,50000)=0 then
			select concat('Var=',Var);
		end if;
		set lend_status = 0;
    set j = floor(rand()*10000+1); -- 1-10000
    if j <= 5 then
			set lend_status = 5; -- 已下架
		end if;
		set cid=floor(rand()*6+6);
        insert into bookms_book(book_name,identify,`description`,catalog,cover,`status`,sort,create_time,doc_count,comment_count,favorite_count,score,score_count,author,book_count,average_score) values("test book", CONCAT('test-identify', Var),"test description", "test catalog",'',0,1,Now(),0,0,0,0,0,"GrassInWind",3,0);
        insert into bookms_book_record(identify,`lend_status`,`user_id`,lend_time,return_time,lend_count,store_position) values(CONCAT('test-identify', Var),lend_status,user_id,Now(),Now(),1, "test store position");
        insert into bookms_book_category(identify, category_id) values(CONCAT('test-identify', Var), cid);
        set k = floor(rand()*100+1);
        if k <= 20 then
			    insert into bookms_user_favorite(`user_id`,identify) values(user_id, CONCAT('test-identify', Var));
        end if;
        set Var = Var + 1;
	end while loop1;
end;
call BatchInsert(1000000);
-- delimiter ;
