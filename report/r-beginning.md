---
figureTitle: Рисунок
tableTitle: Таблица
tableEqns: true
titleDelim: "&nbsp;–"
link-citations: true
linkReferences: true
chapters: true
...


# [Реферат]{custom-style="UnnumberedHeadingOneNoTOC"} {.unnumbered}

Выпускная квалификационная работа состоит из %NPAGES% страниц, %NFIGURES% рисунков, %NTABLES% таблиц, %NREFERENCES% источника и %NAPPENDICES% приложений.

::: {custom-style="GostKeywords"}
миграции,
реляционные базы данных,
микросервисная архитектура,
чистая архитектура,
REST API,
gRPC,
PostgreSQL,
контейнеризация,
автоматизация,
документация
:::

Объект исследования — процесс управления миграциями в реляционных базах данных.

Предмет исследования — серверная компонента приложения, реализующая управление миграциями и обеспечивающая их мониторинг, автоматизацию и интеграцию.

Цель работы — разработка серверного приложения для управления миграциями в реляционных базах данных с использованием микросервисной архитектуры, принципов чистой архитектуры и современных технологий, обеспечивающих автоматизацию, мониторинг и удобство интеграции.

Для достижения цели проанализированы инструменты Migrate и Goose, выявлены их ограничения, сформулированы требования к системе. Спроектирована и разработана серверная компонента на основе микросервисной и чистой архитектуры.

Основной результат работы – микросервисная система для управления миграциями в PostgreSQL через REST API и gRPC, с поддержкой аутентификации, авторизации, логирования и автоматической генерации документации, развернутая в контейнерах.

Результаты разработки предназначены для использования в проектах, требующих автоматизированного управления изменениями структуры баз данных, особенно в условиях современных CI/CD-процессов. Применение системы позволяет повысить надежность, прозрачность и эффективность управления миграциями.

По результатам работы составлена данная работа, состоящая из введения, заключения
и трех глав.

# [Содержание]{custom-style="UnnumberedHeadingOneNoTOC"} {.unnumbered}

%TOC%

# [Перечень сокращений и обозначений]{custom-style="UnnumberedHeadingOneNoTOC"} {.unnumbered}

В настоящей итоговой аттестационной работе применяют следующие сокращения и обозначения:

БД -- База данных

РБД -- Реляционная база данных

СУБД  -- Система управления базами данных

API -- Application Programming Interface (Программный интерфейс приложения)

CLI -- Command Line Interface (Интерфейс командной строки)

CI/CD -- Continuous Integration/Continuous Delivery (Непрерывная интеграция и непрерывная доставка)

JSON  -- JavaScript Object Notation (Нотация объектов JavaScript)

JWT -- JSON Web Token (Веб-токен JSON)

SQL -- Structured Query Language (Структурированный язык запросов)

HTTP  -- Hypertext Transfer Protocol (Протокол передачи гипертекста)

REST  -- Representational State Transfer (Передача репрезентативного состояния)

gRPC  -- google Remote Procedure Calling (Удаленный в  ызов процедур Google)

ACID  -- Atomicity, Consistency, Isolation, Durability (Атомарность, Согласованность, Изоляция, Долговечность)

DDD -- Domain-Driven Design (Предметно-ориентированное проектирование)

UI  -- User Interface (Пользовательский интерфейс)

CRUD  -- Create, Read, Update, Delete (Создание, Чтение, Обновление, Удаление)


# [Введение]{custom-style="UnnumberedHeadingOne"} {.unnumbered}

Разработка программного обеспечения, использующего реляционные базы данных (РБД), требует эффективного управления изменениями их структуры. Миграции, представляющие собой процесс создания, модификации или удаления элементов схемы базы данных, являются неотъемлемой частью жизненного цикла приложений. Они обеспечивают синхронизацию изменений между различными средами — разработки, тестирования и продакшн, поддерживая консистентность данных и упрощая развертывание обновлений.

Отсутствие централизованных инструментов управления миграциями приводит к ряду проблем: несогласованности данных между средами, сложностям с откатом изменений, недостаточной прозрачности операций и ограниченной автоматизации. Существующие решения, такие как Migrate и Goose, обладают недостатками, включая отсутствие готового API для интеграции, минимальные возможности мониторинга состояния миграций и ограниченную поддержку сложных сценариев управления.

Актуальность работы обусловлена необходимостью создания серверной компоненты приложения, которая обеспечит централизованное управление миграциями в реляционных базах данных, предоставит удобный интерфейс для мониторинга их статуса и поддержит автоматизацию процессов. Разработанное решение должно устранять ограничения существующих инструментов, предлагая современный подход к управлению миграциями с учетом требований безопасности, масштабируемости и интеграции.

Цель работы — разработка серверного приложения для управления миграциями в реляционных базах данных с использованием микросервисной архитектуры, принципов чистой архитектуры и современных технологий, обеспечивающих автоматизацию, мониторинг и удобство интеграции.

Для достижения цели поставлены следующие задачи:

- проведение анализа существующих инструментов управления миграциями и выявление их ограничений,
- обоснование выбора микросервисной архитектуры и принципов чистой архитектуры для проектирования системы,
- проектирование структуры базы данных для хранения информации о миграциях и управления доступом,
- реализация микросервисов для управления миграциями и авторизации пользователей,
- разработка REST API и gRPC для взаимодействия с приложением,
- обеспечение контейнеризации приложения с использованием Docker для стандартизации развертывания,
- автоматизация генерации документации API и кода с использованием Swagger и Godoc,
- тестирование и отладка приложения для обеспечения его надежности,
- подготовка документации и инструкций по эксплуатации системы.

Объект исследования — процесс управления миграциями в реляционных базах данных. Предмет исследования — серверная компонента приложения, реализующая управление миграциями и обеспечивающая их мониторинг, автоматизацию и интеграцию.

В работе применяются технологии: язык программирования Go для реализации серверной логики, СУБД PostgreSQL для хранения данных, Docker и Docker Compose для контейнеризации, Swagger и Godoc для автоматической генерации документации. Система построена на основе микросервисной архитектуры, включающей два основных сервиса: сервис управления миграциями и сервис авторизации. Использование принципов чистой архитектуры обеспечивает модульность, независимость компонентов и удобство тестирования. Доступ к операциям с миграциями происходит на основе ролевой модели авторизации с применением JWT-токенов, что гарантирует безопасность и контроль действий. API системы реализовано в форматах REST и gRPC, обеспечивая гибкость интеграции с другими сервисами.

Структура работы включает три основные главы: анализ существующих решений, проектирование системы и описание ее реализации. Каждая глава последовательно раскрывает этапы создания приложения, начиная с анализа аналогов и заканчивая технической реализацией и документацией.

Таким образом, разработка серверной компоненты для управления миграциями представляет собой актуальную задачу, направленную на устранение недостатков существующих решений. Созданное приложение обеспечивает прозрачность, автоматизацию и безопасность процессов управления миграциями, предоставляя удобный API для интеграции и масштабируемую архитектуру для поддержки сложных инфраструктур.