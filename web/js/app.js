// Tick - To Do App
const app = {
    // Текущее состояние приложения
    currentFilter: 'all',
    currentTodo: null,
    todos: [],
    searchQuery: '',
    
    // Инициализация приложения
    init: function() {
        console.log('Tick приложение инициализировано');
        this.loadTodos();
        this.bindEvents();
        this.updateStats();
        
        // Скрываем приветственный экран на мобильных при первой загрузке
        if (window.innerWidth <= 768) {
            this.hideWelcome();
        }
    },
    
    // Привязка обработчиков событий
    bindEvents: function() {
        // Обработчик изменения размера окна
        window.addEventListener('resize', () => {
            this.handleResize();
        });
        
        // Обработчики для клавиш Enter в формах
        document.getElementById('new-title').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                app.createTodo();
            }
        });
        
        document.getElementById('todo-title').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                app.saveTodo();
            }
        });
    },
    
    // Обработка изменения размера окна
    handleResize: function() {
        if (window.innerWidth > 768) {
            // На десктопе показываем оба блока
            document.getElementById('main-area').classList.remove('active');
            document.getElementById('back-btn').style.display = 'none';
            // Показываем приветствие если нет активной задачи
            if (!this.currentTodo && !document.getElementById('todo-create').classList.contains('active')) {
                this.showWelcome();
            }
        } else {
            // На мобильных скрываем приветствие
            this.hideWelcome();
        }
    },
    
    // Переключение боковой панели на мобильных
    toggleSidebar: function() {
        const sidebar = document.getElementById('sidebar');
        if (window.innerWidth <= 768) {
            sidebar.classList.toggle('open');
        }
    },
    
    // Показать детали задачи
    showTodoDetail: function(todo) {
        console.log('Показ деталей задачи:', todo);
        this.currentTodo = todo;
        
        // Заполняем форму данными задачи
        document.getElementById('todo-title').value = todo.title || '';
        document.getElementById('todo-completed').checked = todo.completed || false;
        document.getElementById('todo-description').value = todo.description || '';
        document.getElementById('todo-date').textContent = this.formatDate(todo.created_at);
        
        // Скрываем все блоки
        this.hideWelcome();
        document.getElementById('todo-create').classList.remove('active');
        document.getElementById('todo-detail').classList.remove('active');
        
        // Показываем детали
        document.getElementById('todo-detail').classList.add('active');
        
        // На мобильных устройствах переключаемся на экран деталей
        if (window.innerWidth <= 768) {
            document.getElementById('main-area').classList.add('active');
            document.getElementById('back-btn').style.display = 'block';
            document.getElementById('back-btn').textContent = '← К списку';
            // Закрываем sidebar если он открыт
            document.getElementById('sidebar').classList.remove('open');
        }
        
        // Обновляем активный элемент в списке
        this.updateActiveTodo(todo.id);
    },
    
    // Скрыть детали задачи
    hideTodoDetail: function() {
        document.getElementById('todo-detail').classList.remove('active');
        this.currentTodo = null;
        
        if (window.innerWidth <= 768) {
            document.getElementById('main-area').classList.remove('active');
            document.getElementById('back-btn').style.display = 'none';
        }
        
        this.showWelcome();
        this.updateActiveTodo(null);
    },
    
    // Показать форму создания новой задачи
    showNewTodo: function() {
        console.log('Показ формы создания задачи');
        // Очищаем форму
        document.getElementById('new-title').value = '';
        document.getElementById('new-description').value = '';
        
        // Скрываем все блоки
        this.hideWelcome();
        document.getElementById('todo-detail').classList.remove('active');
        document.getElementById('todo-create').classList.remove('active');
        
        // Показываем форму создания
        document.getElementById('todo-create').classList.add('active');
        
        // На мобильных устройствах переключаемся на экран создания
        if (window.innerWidth <= 768) {
            document.getElementById('main-area').classList.add('active');
            document.getElementById('back-btn').style.display = 'block';
            document.getElementById('back-btn').textContent = '← Отмена';
            // Закрываем sidebar если он открыт
            document.getElementById('sidebar').classList.remove('open');
        }
        
        // Фокусируемся на поле ввода названия
        setTimeout(() => {
            document.getElementById('new-title').focus();
        }, 100);
    },
    
    // Скрыть приветственный экран
    hideWelcome: function() {
        document.getElementById('welcome').style.display = 'none';
    },
    
    // Показать приветственный экран
    showWelcome: function() {
        if (window.innerWidth > 768) {
            document.getElementById('welcome').style.display = 'flex';
        }
    },
    
    // Обновить активную задачу в списке
    updateActiveTodo: function(todoId) {
        // Убираем активный класс у всех задач
        document.querySelectorAll('.todo-item').forEach(item => {
            item.classList.remove('active');
        });
        
        // Добавляем активный класс выбранной задаче
        if (todoId) {
            const activeItem = document.querySelector(`.todo-item[data-id="${todoId}"]`);
            if (activeItem) {
                activeItem.classList.add('active');
            }
        }
    },
    
    // Отменить редактирование
    cancelEdit: function() {
        this.hideTodoDetail();
    },
    
    // Отменить создание
    cancelCreate: function() {
        document.getElementById('todo-create').classList.remove('active');
        
        if (window.innerWidth <= 768) {
            document.getElementById('main-area').classList.remove('active');
            document.getElementById('back-btn').style.display = 'none';
        }
        
        this.showWelcome();
    },
    
    // Загрузить задачи с сервера
    loadTodos: function() {
        console.log('Загрузка задач...');
        
        fetch('/api/todos', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Задачи загружены:', data);
            this.todos = Array.isArray(data) ? data : [];
            this.renderTodos();
            this.updateStats();
            
            // Если есть задачи и на десктопе, показываем приветствие
            if (this.todos.length === 0 && window.innerWidth > 768) {
                this.showWelcome();
            }
        })
        .catch(error => {
            console.error('Ошибка загрузки задач:', error);
            this.showToast('Ошибка загрузки задач', 'error');
            
            // Показываем состояние ошибки
            const todoList = document.getElementById('todo-list');
            todoList.innerHTML = `
                <div class="empty-state">
                    <div class="empty-icon">⚠️</div>
                    <div class="empty-title">Ошибка загрузки</div>
                    <div class="empty-description">Не удалось загрузить задачи. Проверьте подключение к интернету.</div>
                    <button class="primary-btn" onclick="app.loadTodos()">⟳ Повторить</button>
                </div>
            `;
        });
    },
    
    // Создать новую задачу
    createTodo: function() {
        const title = document.getElementById('new-title').value.trim();
        const description = document.getElementById('new-description').value.trim();
        
        if (!title) {
            this.showToast('Введите название задачи', 'warning');
            document.getElementById('new-title').focus();
            return;
        }
        
        const todoData = {
            title: title,
            description: description,
            completed: false
        };
        
        console.log('Создание задачи:', todoData);
        
        fetch('/api/todos', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(todoData)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Задача создана:', data);
            this.showToast('Задача создана успешно', 'success');
            this.loadTodos();
            this.cancelCreate();
        })
        .catch(error => {
            console.error('Ошибка создания задачи:', error);
            this.showToast('Ошибка создания задачи', 'error');
        });
    },
    
    // Сохранить изменения в задаче
    saveTodo: function() {
        if (!this.currentTodo) {
            this.showToast('Нет активной задачи для сохранения', 'warning');
            return;
        }
        
        const title = document.getElementById('todo-title').value.trim();
        const description = document.getElementById('todo-description').value.trim();
        const completed = document.getElementById('todo-completed').checked;
        
        if (!title) {
            this.showToast('Введите название задачи', 'warning');
            document.getElementById('todo-title').focus();
            return;
        }
        
        const todoData = {
            title: title,
            description: description,
            completed: completed
        };
        
        console.log('Сохранение задачи:', todoData);
        
        fetch(`/api/todos/${this.currentTodo.id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(todoData)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Задача обновлена:', data);
            this.showToast('Изменения сохранены', 'success');
            this.loadTodos();
            
            // Обновляем текущую задачу
            if (data && data.id) {
                this.currentTodo = data;
            }
        })
        .catch(error => {
            console.error('Ошибка сохранения задачи:', error);
            this.showToast('Ошибка сохранения задачи', 'error');
        });
    },
    
    // Удалить задачу
    deleteTodo: function() {
        if (!this.currentTodo) {
            this.showToast('Нет активной задачи для удаления', 'warning');
            return;
        }
        
        if (!confirm('Вы уверены, что хотите удалить эту задачу?')) {
            return;
        }
        
        console.log('Удаление задачи:', this.currentTodo.id);
        
        fetch(`/api/todos/${this.currentTodo.id}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.text();
        })
        .then(() => {
            console.log('Задача удалена');
            this.showToast('Задача удалена', 'success');
            this.hideTodoDetail();
            this.loadTodos();
        })
        .catch(error => {
            console.error('Ошибка удаления задачи:', error);
            this.showToast('Ошибка удаления задачи', 'error');
        });
    },
    
    // Переключить статус выполнения задачи
    toggleTodoComplete: function(todoId) {
        event.stopPropagation();
        
        const todo = this.todos.find(t => t.id === todoId);
        if (!todo) return;
        
        const newCompletedState = !todo.completed;
        
        console.log('Переключение статуса задачи:', todoId, newCompletedState);
        
        fetch(`/api/todos/${todoId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                ...todo,
                completed: newCompletedState
            })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Статус обновлен:', data);
            this.loadTodos();
            
            // Если это текущая задача, обновляем UI
            if (this.currentTodo && this.currentTodo.id === todoId) {
                this.currentTodo.completed = newCompletedState;
                document.getElementById('todo-completed').checked = newCompletedState;
            }
        })
        .catch(error => {
            console.error('Ошибка обновления статуса:', error);
            this.showToast('Ошибка обновления статуса', 'error');
        });
    },
    
    // Фильтрация задач по поисковому запросу
    filterTodos: function() {
        this.searchQuery = document.getElementById('searchInput').value.toLowerCase().trim();
        this.renderTodos();
        this.updateStats();
    },
    
    // Установить фильтр (все/активные/завершенные)
    setFilter: function(filter) {
        console.log('Установка фильтра:', filter);
        this.currentFilter = filter;
        
        // Обновляем активные кнопки фильтров
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        event.target.classList.add('active');
        
        this.renderTodos();
        this.updateStats();
    },
    
    // Получить отфильтрованные задачи
    getFilteredTodos: function() {
        let filtered = this.todos;
        
        // Применяем текстовый поиск
        if (this.searchQuery) {
            filtered = filtered.filter(todo => 
                todo.title.toLowerCase().includes(this.searchQuery) ||
                (todo.description && todo.description.toLowerCase().includes(this.searchQuery))
            );
        }
        
        // Применяем фильтр по статусу
        switch (this.currentFilter) {
            case 'active':
                filtered = filtered.filter(todo => !todo.completed);
                break;
            case 'completed':
                filtered = filtered.filter(todo => todo.completed);
                break;
            // 'all' - не фильтруем по статусу
        }
        
        return filtered;
    },
    
    // Обновить статистику
    updateStats: function() {
        const total = this.todos.length;
        const active = this.todos.filter(todo => !todo.completed).length;
        const filtered = this.getFilteredTodos();
        
        document.getElementById('total-count').textContent = filtered.length;
        document.getElementById('active-count').textContent = filtered.filter(todo => !todo.completed).length;
        
        // Обновляем заголовок на мобильных
        if (window.innerWidth <= 768) {
            const headerTitle = document.querySelector('.sidebar-header h3');
            if (headerTitle) {
                headerTitle.textContent = `Мои задачи (${filtered.length})`;
            }
        }
    },
    
    // Отобразить задачи
    renderTodos: function() {
        const todoList = document.getElementById('todo-list');
        const filteredTodos = this.getFilteredTodos();
        
        if (filteredTodos.length === 0) {
            let message = '';
            let showButton = false;
            
            if (this.searchQuery) {
                message = 'По вашему запросу ничего не найдено';
                showButton = false;
            } else {
                switch (this.currentFilter) {
                    case 'all':
                        message = 'У вас пока нет задач';
                        showButton = true;
                        break;
                    case 'active':
                        message = 'Все задачи завершены!';
                        showButton = false;
                        break;
                    case 'completed':
                        message = 'Нет завершенных задач';
                        showButton = false;
                        break;
                }
            }
            
            const emptyState = `
                <div class="empty-state">
                    <div class="empty-icon">${showButton ? '📝' : '🔍'}</div>
                    <div class="empty-title">${message}</div>
                    ${showButton ? `
                        <button class="primary-btn" onclick="app.showNewTodo()">
                            + Создать первую задачу
                        </button>
                    ` : ''}
                </div>
            `;
            todoList.innerHTML = emptyState;
            return;
        }
        
        let html = '';
        filteredTodos.forEach(todo => {
            const isActive = this.currentTodo && this.currentTodo.id === todo.id;
            const dateStr = this.formatDate(todo.created_at);
            const descriptionPreview = todo.description ? 
                this.escapeHtml(todo.description.substring(0, 100)) + (todo.description.length > 100 ? '...' : '') : 
                '';
            
            html += `
                <div class="todo-item ${todo.completed ? 'completed' : ''} ${isActive ? 'active' : ''}" 
                     data-id="${todo.id}" onclick="app.showTodoDetail(${JSON.stringify(todo).replace(/"/g, '&quot;')})">
                    <div class="todo-item-header">
                        <div class="todo-checkbox ${todo.completed ? 'checked' : ''}" 
                             onclick="event.stopPropagation(); app.toggleTodoComplete(${todo.id})"
                             title="${todo.completed ? 'Отметить как активную' : 'Отметить как выполненную'}"></div>
                        <div class="todo-item-title" title="${this.escapeHtml(todo.title)}">
                            ${this.escapeHtml(todo.title)}
                        </div>
                    </div>
                    ${descriptionPreview ? `
                        <div class="todo-description-preview" title="${this.escapeHtml(todo.description)}">
                            ${descriptionPreview}
                        </div>
                    ` : ''}
                    <div class="todo-item-date" title="Создано ${dateStr}">
                        ${dateStr}
                    </div>
                </div>
            `;
        });
        
        todoList.innerHTML = html;
    },
    
    // Показать уведомление
    showToast: function(message, type = 'info') {
        const toast = document.getElementById('toast');
        const toastBody = document.getElementById('toast-body');
        
        // Устанавливаем сообщение
        toastBody.textContent = message;
        
        // Устанавливаем цвет в зависимости от типа
        const typeColors = {
            success: '#10b981',
            error: '#ef4444',
            warning: '#f59e0b',
            info: '#3b82f6'
        };
        
        const header = toast.querySelector('.toast-header');
        if (header) {
            header.style.borderLeft = `4px solid ${typeColors[type] || typeColors.info}`;
        }
        
        // Показываем toast
        toast.classList.add('show');
        
        // Автоматически скрываем через 5 секунд
        setTimeout(() => {
            this.hideToast();
        }, 5000);
    },
    
    // Скрыть уведомление
    hideToast: function() {
        const toast = document.getElementById('toast');
        toast.classList.remove('show');
    },
    
    // Форматирование даты
    formatDate: function(dateString) {
        if (!dateString) return '';
        
        const date = new Date(dateString);
        if (isNaN(date.getTime())) return '';
        
        const now = new Date();
        const diffMs = now - date;
        const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
        
        // Если сегодня
        if (diffDays === 0) {
            const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
            if (diffHours < 1) {
                const diffMinutes = Math.floor(diffMs / (1000 * 60));
                return diffMinutes < 1 ? 'только что' : `${diffMinutes} мин назад`;
            }
            return `${diffHours} ч назад`;
        }
        
        // Если вчера
        if (diffDays === 1) {
            return 'вчера';
        }
        
        // Если на этой неделе
        if (diffDays < 7) {
            return `${diffDays} дн назад`;
        }
        
        // Более недели назад - показываем дату
        const options = { 
            day: 'numeric', 
            month: 'short',
            year: now.getFullYear() !== date.getFullYear() ? 'numeric' : undefined
        };
        return date.toLocaleDateString('ru-RU', options);
    },
    
    // Экранирование HTML
    escapeHtml: function(text) {
        if (!text) return '';
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
};

// Инициализация приложения при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    app.init();
    
    // Обработчик клика вне sidebar на мобильных
    document.addEventListener('click', (e) => {
        const sidebar = document.getElementById('sidebar');
        const menuButton = document.querySelector('.menu-button');
        
        if (window.innerWidth <= 768 && 
            sidebar.classList.contains('open') && 
            !sidebar.contains(e.target) && 
            !menuButton.contains(e.target)) {
            sidebar.classList.remove('open');
        }
    });
    
    // Закрытие sidebar при выборе задачи на мобильных
    document.addEventListener('click', (e) => {
        if (window.innerWidth <= 768 && e.target.closest('.todo-item')) {
            const sidebar = document.getElementById('sidebar');
            sidebar.classList.remove('open');
        }
    });
});

// Глобальные функции для доступа из HTML
window.app = app;