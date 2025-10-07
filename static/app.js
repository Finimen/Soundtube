const API_BASE = '/api';
let currentToken = localStorage.getItem('authToken');
let currentUserName = localStorage.getItem('userName');
let currentSoundId = null;

document.addEventListener('DOMContentLoaded', function() {
    console.log('App loaded, token exists:', !!currentToken);
    checkAuth();
    
    if (currentToken) {
        loadSounds();
    } else {
        showUnauthorizedMessage();
    }
    
    document.getElementById('authModal').addEventListener('click', function(e) {
        if (e.target === this) {
            hideAuthModal();
        }
    });
});

function showUnauthorizedMessage() {
    const soundsList = document.getElementById('soundsList');
    soundsList.innerHTML = `
        <h3>Последние треки</h3>
        <div style="text-align: center; padding: 40px;">
            <p>Для просмотра треков необходимо авторизоваться</p>
            <button class="btn btn-primary" onclick="showAuthModal('login')">Войти</button>
            <button class="btn btn-secondary" onclick="showAuthModal('register')">Зарегистрироваться</button>
        </div>
    `;
}

function checkAuth() {
    if (currentToken && currentUserName) {
        document.getElementById('authSection').classList.add('hidden');
        document.getElementById('userInfo').classList.remove('hidden');
        document.getElementById('uploadSection').classList.remove('hidden');
        document.getElementById('userGreeting').textContent = `Привет, ${currentUserName}!`;
    } else {
        document.getElementById('authSection').classList.remove('hidden');
        document.getElementById('userInfo').classList.add('hidden');
        document.getElementById('uploadSection').classList.add('hidden');
    }
}

function showAuthModal(type) {
    const modal = document.getElementById('authModal');
    const title = document.getElementById('modalTitle');
    const form = document.getElementById('authForm');
    const nameField = document.getElementById('nameField');
    const emailField = document.getElementById('emailField');
    const usernameField = document.getElementById('usernameField');
    const passwordNote = document.getElementById('passwordNote');
    
    form.reset();
    
    if (type === 'register') {
        title.textContent = 'Регистрация';
        nameField.classList.remove('hidden');
        emailField.classList.remove('hidden');
        usernameField.classList.add('hidden');
        passwordNote.textContent = 'Минимум 6 символов';
    } else {
        title.textContent = 'Вход';
        nameField.classList.add('hidden');
        emailField.classList.add('hidden');
        usernameField.classList.remove('hidden');
        passwordNote.textContent = 'Введите ваш пароль';
    }
    
    modal.classList.remove('hidden');
    
    form.onsubmit = function(e) {
        e.preventDefault();
        if (type === 'register') {
            register();
        } else {
            login();
        }
    };
}

function hideAuthModal() {
    const modal = document.getElementById('authModal');
    modal.classList.add('hidden');
    document.getElementById('authForm').reset();
}

async function register() {
    const username = document.getElementById('registerName').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    console.log('Registration data:', { username, email, password }); 

    if (username.length < 3 || username.length > 50) {
        alert('Имя пользователя должно быть от 3 до 50 символов');
        return;
    }

    if (password.length < 6) {
        alert('Пароль должен содержать минимум 6 символов');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                username: username, 
                email: email, 
                password: password 
            })
        });

        if (response.ok) {
            alert('Регистрация успешна! Теперь войдите в систему.');
            hideAuthModal();
            setTimeout(() => showAuthModal('login'), 500);
        } else {
            const errorData = await response.json();
            alert('Ошибка регистрации: ' + (errorData.error || 'Неизвестная ошибка'));
        }
    } catch (error) {
        alert('Ошибка сети: ' + error.message);
    }
}

async function login() {
    const username = document.getElementById('loginUsername').value; 
    const password = document.getElementById('password').value;

    try {
        const response = await fetch(`${API_BASE}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                username: username,
                password: password 
            })
        });

        if (response.ok) {
            const data = await response.json();
            console.log('Login response:', data); 
            
            currentToken = data.token || data; 
            currentUserName = username; 
            
            localStorage.setItem('authToken', currentToken);
            localStorage.setItem('userName', currentUserName);
            
            console.log('Token saved:', currentToken);
            
            hideAuthModal();
            checkAuth();
            loadSounds();
        } else {
            const error = await response.text();
            alert('Ошибка входа: ' + error);
        }
    } catch (error) {
        alert('Ошибка сети: ' + error.message);
    }
}

async function logout() {
    try {
        await fetch(`${API_BASE}/auth/logout`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });
    } catch (error) {
        console.error('Ошибка выхода:', error);
    }
    
    currentToken = null;
    currentUserName = null;
    localStorage.removeItem('authToken');
    localStorage.removeItem('userName');
    checkAuth();
    loadSounds();
}

async function loadSounds() {
    const soundsList = document.getElementById('soundsList');
    soundsList.innerHTML = '<h3>Последние треки</h3>';

    try {
        const headers = {};
        if (currentToken) {
            headers['Authorization'] = `Bearer ${currentToken}`;
        }

        const response = await fetch(`${API_BASE}/sounds/`, { headers });
        
        if (response.ok) {
            const sounds = await response.json();
            if (sounds.length === 0) {
                soundsList.innerHTML += '<p>Пока нет загруженных треков. Будьте первым!</p>';
            } else {
                sounds.forEach(sound => {
                    const soundElement = createSoundElement(sound);
                    soundsList.appendChild(soundElement);
                });
            }
        } else if (response.status === 401) {
            soundsList.innerHTML += '<p>Для просмотра треков необходимо авторизоваться</p>';
        }
    } catch (error) {
        console.error('Ошибка загрузки звуков:', error);
        soundsList.innerHTML += '<p>Ошибка загрузки треков</p>';
    }
}

function createSoundElement(sound) {
    const div = document.createElement('div');
    div.className = 'sound-item';
    div.innerHTML = `
        <div class="sound-title">${escapeHtml(sound.title)}</div>
        <div class="sound-description">${escapeHtml(sound.description || '')}</div>
        <div class="sound-author" style="color: #888; font-size: 0.9rem; margin-bottom: 10px;">
            Автор: ${escapeHtml(sound.author_name || 'Неизвестен')}
        </div>
        <audio controls style="width: 100%; margin: 10px 0;">
            <source src="/static/${sound.file_path}" type="audio/mpeg">
            Ваш браузер не поддерживает аудио элементы.
        </audio>
        <div class="comments-section">
            <h4>Комментарии</h4>
            <div id="comments-${sound.id}">
                <!-- Комментарии будут загружены здесь -->
            </div>
            ${currentToken ? `
                <div class="comment-form">
                    <textarea class="form-control" id="comment-text-${sound.id}" placeholder="Добавить комментарий..."></textarea>
                    <button class="btn btn-primary" onclick="addComment(${sound.id})" style="margin-top: 10px;">Отправить</button>
                </div>
            ` : '<p>Войдите, чтобы комментировать</p>'}
        </div>
    `;
    
    loadComments(sound.id);
    
    return div;
}

async function loadComments(soundId) {
    if (!currentToken) return;

    try {
        const response = await fetch(`${API_BASE}/sounds/${soundId}/comments`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        if (response.ok) {
            const comments = await response.json();
            const commentsContainer = document.getElementById(`comments-${soundId}`);
            commentsContainer.innerHTML = '';
            
            if (comments.length === 0) {
                commentsContainer.innerHTML = '<p>Пока нет комментариев. Будьте первым!</p>';
            } else {
                comments.forEach(comment => {
                    const commentElement = document.createElement('div');
                    commentElement.className = 'comment';
                    commentElement.innerHTML = `
                        <strong>${escapeHtml(comment.author_name)}</strong>
                        <p>${escapeHtml(comment.text)}</p>
                    `;
                    commentsContainer.appendChild(commentElement);
                });
            }
        }
    } catch (error) {
        console.error('Ошибка загрузки комментариев:', error);
    }
}

async function addComment(soundId) {
    if (!currentToken) {
        alert('Для комментирования необходимо авторизоваться');
        return;
    }

    const text = document.getElementById(`comment-text-${soundId}`).value;
    if (!text.trim()) {
        alert('Введите текст комментария');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/sounds/${soundId}/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${currentToken}`
            },
            body: JSON.stringify({ text })
        });

        if (response.ok) {
            document.getElementById(`comment-text-${soundId}`).value = '';
            loadComments(soundId);
        } else {
            alert('Ошибка при добавлении комментария');
        }
    } catch (error) {
        alert('Ошибка сети: ' + error.message);
    }
}

document.getElementById('uploadForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    if (!currentToken) {
        alert('Для загрузки трека необходимо авторизоваться');
        return;
    }

    const formData = new FormData();
    formData.append('title', document.getElementById('soundTitle').value);
    formData.append('description', document.getElementById('soundDescription').value);
    formData.append('file', document.getElementById('soundFile').files[0]);

    try {
        const response = await fetch(`${API_BASE}/sounds/`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${currentToken}`
            },
            body: formData
        });

        if (response.ok) {
            alert('Трек успешно загружен!');
            document.getElementById('uploadForm').reset();
            loadSounds();
        } else {
            const error = await response.text();
            alert('Ошибка загрузки: ' + error);
        }
    } catch (error) {
        alert('Ошибка сети: ' + error.message);
    }
});

function showUploadSection() {
    if (!currentToken) {
        showAuthModal('login');
        return;
    }
    document.getElementById('uploadSection').scrollIntoView({ behavior: 'smooth' });
}

function escapeHtml(unsafe) {
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}