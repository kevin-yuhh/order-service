INSERT INTO strategy (id, type, lua_script)
VALUES (1, 0, 'M = 1000000  -- 1 MB
G = M * 1000 -- 1 GB
D = 5        -- upload times > 1000, file size <= 1M charge 0.000005BTT
F = 1998     -- 1 GB charge
T = 1000     -- 1000 times

function size_fee(size, total_times)
    local default_charge = 0
    if total_times > T then
        default_charge = D
    end

    if size > 0 and size <= M then
        return default_charge;
    elseif size > M and size <= G then
        return (size - M) / M * 2 + default_charge;
    else
        return default_charge + F + (size - G) / M;
    end
end

function calc(size, total_times, time)
    return time*(size_fee(size, total_times));
end');

INSERT INTO config (env, strategy_id, default_time) VALUES ('TEST',1,30);