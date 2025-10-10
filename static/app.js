const API_BASE = '/api';
let currentToken = localStorage.getItem('authToken');
let currentUserName = localStorage.getItem('userName');
let currentSoundId = null;
let currentCommentsSoundId = null;

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

    document.getElementById('commentsModal').addEventListener('click', function(e) {
        if (e.target === this) {
            hideCommentsModal();
        }
    });

    // Initialize drag and drop
    initializeFileDrop();
});

function initializeFileDrop() {
    const dropArea = document.getElementById('fileDropArea');
    const fileInput = document.getElementById('soundFile');

    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        dropArea.addEventListener(eventName, preventDefaults, false);
    });

    function preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    ['dragenter', 'dragover'].forEach(eventName => {
        dropArea.addEventListener(eventName, highlight, false);
    });

    ['dragleave', 'drop'].forEach(eventName => {
        dropArea.addEventListener(eventName, unhighlight, false);
    });

    function highlight() {
        dropArea.classList.add('highlight');
    }

    function unhighlight() {
        dropArea.classList.remove('highlight');
    }

    dropArea.addEventListener('drop', handleDrop, false);
    dropArea.addEventListener('click', () => fileInput.click());

    fileInput.addEventListener('change', handleFileSelect);

    function handleDrop(e) {
        const dt = e.dataTransfer;
        const files = dt.files;
        fileInput.files = files;
        handleFiles(files);
    }

    function handleFileSelect() {
        handleFiles(this.files);
    }

    function handleFiles(files) {
        if (files.length > 0) {
            const file = files[0];
            if (file.type.startsWith('audio/')) {
                dropArea.querySelector('.file-msg').textContent = `Выбран файл: ${file.name}`;
            } else {
                alert('Пожалуйста, выберите аудио файл');
                fileInput.value = '';
                dropArea.querySelector('.file-msg').textContent = 'Перетащите аудио файл сюда или нажмите для выбора';
            }
        }
    }
}

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
    div.id = `sound-${sound.id}`;

    console.log('Sound data:', sound);

    const soundName = sound.name || sound.title || 'Без названия';
    const soundAlbum = sound.album || 'Не указан';
    const soundGenre = sound.genre || 'Не указан';
    const authorId = sound.author_id || sound.authorID || 'Неизвестен';
    const filePath = sound.file_path || sound.filePath || sound.filename;
    const likes = sound.likes || 0;
    const dislikes = sound.dislikes || 0;
    const userReaction = sound.user_reaction || null;

    div.innerHTML = `
        <div class="sound-title">${escapeHtml(soundName)}</div>
        <div class="sound-info">
            <span class="sound-album">Альбом: ${escapeHtml(soundAlbum)}</span>
            <span class="sound-genre">Жанр: ${escapeHtml(soundGenre)}</span>
        </div>
        <div class="sound-author" style="color: #888; font-size: 0.9rem; margin-bottom: 10px;">
            Автор ID: ${authorId}
        </div>
        ${filePath ? `
            <audio controls style="width: 100%; margin: 10px 0;">
                <source src="/static/${filePath}" type="audio/mpeg">
                Ваш браузер не поддерживает аудио элементы.
            </audio>
        ` : '<p>Аудио файл не загружен</p>'}
        
        <div class="sound-actions">
            <div class="reactions">
                <button class="reaction-btn ${userReaction === 'like' ? 'active' : ''}" onclick="setReaction(${sound.id}, 'like')">
                    👍 ${likes}
                </button>
                <button class="reaction-btn ${userReaction === 'dislike' ? 'active' : ''}" onclick="setReaction(${sound.id}, 'dislike')">
                    👎 ${dislikes}
                </button>
                <button class="comment-btn" onclick="showComments(${sound.id}, '${escapeHtml(soundName)}')">
                    💬 Комментарии
                </button>
            </div>
        </div>
    `;

    return div;
}

// Функции для реакций
async function setReaction(soundId, reactionType) {
    if (!currentToken) {
        alert('Для оценки трека необходимо авторизоваться');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/sounds/${soundId}/reactions`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${currentToken}`
            },
            body: JSON.stringify({ type: reactionType })
        });

        if (response.ok) {
            // Обновляем отображение реакций
            await updateSoundReactions(soundId);
        } else {
            const error = await response.text();
            alert('Ошибка установки реакции: ' + error);
        }
    } catch (error) {
        alert('Ошибка сети: ' + error.message);
    }
}

async function deleteReaction(soundId) {
    if (!currentToken) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/sounds/${soundId}/reactions`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        if (response.ok) {
            await updateSoundReactions(soundId);
        }
    } catch (error) {
        console.error('Ошибка удаления реакции:', error);
    }
}

async function updateSoundReactions(soundId) {
    try {
        const response = await fetch(`${API_BASE}/sounds/${soundId}/reactions`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        if (response.ok) {
            const reactions = await response.json();
            const soundElement = document.getElementById(`sound-${soundId}`);
            if (soundElement) {
                const likeBtn = soundElement.querySelector('.reaction-btn:nth-child(1)');
                const dislikeBtn = soundElement.querySelector('.reaction-btn:nth-child(2)');

                likeBtn.textContent = `👍 ${reactions.likes || 0}`;
                dislikeBtn.textContent = `👎 ${reactions.dislikes || 0}`;

                likeBtn.classList.toggle('active', reactions.user_reaction === 'like');
                dislikeBtn.classList.toggle('active', reactions.user_reaction === 'dislike');
            }
        }
    } catch (error) {
        console.error('Ошибка обновления реакций:', error);
    }
}

// Функции для комментариев
async function showComments(soundId, soundName) {
    if (!currentToken) {
        alert('Для просмотра комментариев необходимо авторизоваться');
        return;
    }

    currentCommentsSoundId = soundId;
    const modal = document.getElementById('commentsModal');
    const title = document.getElementById('commentsModalTitle');

    title.textContent = `Комментарии: ${soundName}`;
    modal.classList.remove('hidden');

    await loadComments(soundId);
}

function hideCommentsModal() {
    const modal = document.getElementById('commentsModal');
    modal.classList.add('hidden');
    document.getElementById('newCommentText').value = '';
    currentCommentsSoundId = null;
}

async function loadComments(soundId) {
    const commentsList = document.getElementById('commentsList');
    commentsList.innerHTML = '<p>Загрузка комментариев...</p>';

    try {
        const response = await fetch(`${API_BASE}/sounds/${soundId}/comments`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        if (response.ok) {
            const comments = await response.json();
            displayComments(comments);
        } else {
            commentsList.innerHTML = '<p>Ошибка загрузки комментариев</p>';
        }
    } catch (error) {
        console.error('Ошибка загрузки комментариев:', error);
        commentsList.innerHTML = '<p>Ошибка загрузки комментариев</p>';
    }
}

function displayComments(comments) {
    const commentsList = document.getElementById('commentsList');

    if (comments.length === 0) {
        commentsList.innerHTML = '<p>Пока нет комментариев. Будьте первым!</p>';
        return;
    }

    commentsList.innerHTML = comments.map(comment => `
        <div class="comment-item">
            <div class="comment-header">
                <strong>${escapeHtml(comment.author_name || 'Аноним')}</strong>
                <span class="comment-date">${new Date(comment.created_at).toLocaleString()}</span>
            </div>
            <div class="comment-text">${escapeHtml(comment.text)}</div>
            <div class="comment-actions">
                <button class="reaction-btn" onclick="setCommentReaction(${comment.id}, 'like')">
                    👍 ${comment.likes || 0}
                </button>
                <button class="reaction-btn" onclick="setCommentReaction(${comment.id}, 'dislike')">
                    👎 ${comment.dislikes || 0}
                </button>
            </div>
        </div>
    `).join('');
}

async function addComment() {
    if (!currentToken) {
        return;
    }

    const text = document.getElementById('newCommentText').value.trim();
    if (!text) {
        alert('Введите текст комментария');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/sounds/${currentCommentsSoundId}/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${currentToken}`
            },
            body: JSON.stringify({ text: text })
        });

        if (response.ok) {
            document.getElementById('newCommentText').value = '';
            await loadComments(currentCommentsSoundId);
        } else {
            const error = await response.text();
            alert('Ошибка добавления комментария: ' + error);
        }
    } catch (error) {
        alert('Ошибка сети: ' + error.message);
    }
}

// Реакции для комментариев
async function setCommentReaction(commentId, reactionType) {
    if (!currentToken) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/comments/${commentId}/reactions`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${currentToken}`
            },
            body: JSON.stringify({ type: reactionType })
        });

        if (response.ok) {
            await loadComments(currentCommentsSoundId);
        }
    } catch (error) {
        console.error('Ошибка установки реакции на комментарий:', error);
    }
}

document.getElementById('uploadForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    if (!currentToken) {
        alert('Для создания трека необходимо авторизоваться');
        return;
    }

    const name = document.getElementById('soundName').value;
    const album = document.getElementById('soundAlbum').value;
    const genre = document.getElementById('soundGenre').value;
    const file = document.getElementById('soundFile').files[0];

    if (!name || !album || !genre) {
        alert('Пожалуйста, заполните все поля');
        return;
    }

    if (!file) {
        alert('Пожалуйста, выберите аудио файл');
        return;
    }

    try {
        const soundResponse = await fetch(`${API_BASE}/sounds/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${currentToken}`
            },
            body: JSON.stringify({
                name: name,
                album: album,
                genre: genre
            })
        });

        if (soundResponse.ok) {
            const formData = new FormData();
            formData.append('file', file);
            formData.append('name', name);

            const fileResponse = await fetch(`${API_BASE}/sounds/upload`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${currentToken}`
                },
                body: formData
            });

            if (fileResponse.ok) {
                alert('Трек успешно создан и файл загружен!');
                document.getElementById('uploadForm').reset();
                document.getElementById('fileDropArea').querySelector('.file-msg').textContent =
                    'Перетащите аудио файл сюда или нажмите для выбора';
                loadSounds();
            } else {
                alert('Ошибка загрузки файла');
            }
        } else {
            const error = await soundResponse.text();
            alert('Ошибка создания трека: ' + error);
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
    if (!unsafe) return '';
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}