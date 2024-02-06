Требование

Сделать веб сервер для голосования с двумя ручками:
- Для сохранения результатов голосования (vote), принимает номер паспорта и ID кандидата
- Для получения результатов (как для отдельного кандидата, так и по всем кандидатам)

`bombardier -c 100 -d 60s -r 100000 -l -m POST -b '{"passport":"pass","candidate_id": 1}' localhost:8000/vote`
`bombardier -c 125 -n 100000 -l -m POST -b '{"passport":"pass","candidate_id": 1}' localhost:8000/vote`
`bombardier -c 125 -n 100000 -l -m GET localhost:8000/stats`