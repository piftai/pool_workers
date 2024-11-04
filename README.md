# Тестовое задание VK

Реализация примитивного Worker Pool с возможностью динамического добавления/удаления воркеров.

# Использование
Структура WorkerPool имеет атрибуты mu для синхронизации, taskChan как канал с задачами, workers как мапа воркеров.
Хочется отметить, что в первой версии workers был слайсом, но в итоге оптимизировал этот атрибут.
Структура Worker имеет атрибуты id, active, quitChan.


func (p *WorkerPool) StartWorker() - запускает воркера

func (p *WorkerPool) StopWorker(id int) - останавливает воркера с переданным id

func (p *WorkerPool) SetWorkersCount(n int) - для первоначального запуска воркер пула с n воркерами, либо для удаления лишних воркеров

func (p *WorkerPool) AddTask(task string) - для добавления задач в канал

# Работа main.go

В качестве примера в воркер пул загружаем обычные строки, для удобства просмотра добавил time.sleep в некоторых моментах.
Можно посмотреть как динамически добавляется/удаляется воркер.

# Установка и запуск
```bash
    git clone https://github.com/piftai/pool_workers.git
    cd pool_workers
    go run main.go
```
