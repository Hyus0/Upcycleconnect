<template>
    <div class="layout-wrapper public-dashboard">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">MODÉRATION > FORUM</p>
                <h1 class="hero-title1">Modération de la communauté</h1>
                <p class="classic-text">
                    Gérez les signalements, les sujets et l'accès au forum d'UpcycleConnect.
                </p>
            </div>
            <button class="btn-main-action" @click="refreshData" :disabled="isRefreshing">
                {{ isRefreshing ? 'Actualisation...' : 'Actualiser les données' }}
            </button>
        </header>

        <div class="stats-grid dashboard-stats-mod">
            <div class="card card--white" :class="{ 'card-danger-alert': reportedMessages.length > 0 }">
                <div class="card-num" :style="{ color: reportedMessages.length > 0 ? '#d32f2f' : '#1a1a1a' }">
                    {{ reportedMessages.length }}
                </div>
                <p class="text-dm">Alertes</p>
                <span class="badge" :class="reportedMessages.length > 0 ? 'badge--red' : 'badge--green'">
                    {{ reportedMessages.length > 0 ? 'Action requise' : 'Tout est calme' }}
                </span>
            </div>

            <div class="card card--white">
                <div class="card-num">{{ allTopics.length }}</div>
                <p class="text-dm">Discussions</p>
                <span class="badge badge--green">Sujets ouverts</span>
            </div>

            <div class="card card--white">
                <div class="card-num">{{ recentMessages.length }}</div>
                <p class="text-dm">Derniers Msgs</p>
                <span class="badge badge--green">Flux récent</span>
            </div>

            <div class="card card--white">
                <div class="card-num" style="color: #666;">{{ bannedUsers.length }}</div>
                <p class="text-dm">Bannis</p>
                <span class="badge badge--orange">Accès restreints</span>
            </div>
        </div>

        <div class="section-container">
            <div class="section-header">
                <h2>Gestion du contenu</h2>
                <div class="header-actions tabs-actions">
                    <button class="btn-tab" :class="{ 'active': activeTab === 'signalements' }" @click="activeTab = 'signalements'">Alertes</button>
                    <button class="btn-tab" :class="{ 'active': activeTab === 'topics' }" @click="activeTab = 'topics'">Discussions</button>
                    <button class="btn-tab" :class="{ 'active': activeTab === 'messages' }" @click="activeTab = 'messages'">Flux</button>
                    <button class="btn-tab" :class="{ 'active': activeTab === 'bannis' }" @click="activeTab = 'bannis'">Bannis</button>
                </div>
            </div>

            <div v-if="loading" class="state-card" style="margin-top: 1rem;">
                Chargement des données en cours...
            </div>

            <div v-else-if="activeTab === 'signalements'">
                <table v-if="reportedMessages.length > 0" class="data-table">
                    <thead>
                        <tr>
                            <th>AUTEUR & DATE</th>
                            <th>MESSAGE SIGNALÉ</th>
                            <th>MOTIFS</th>
                            <th>ACTIONS</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="msg in reportedMessages" :key="msg.message_id">
                            <td class="user-link" @click="viewProfile(msg.author_id)">
                                <strong>{{ msg.author_name }}</strong><br />
                                <small class="table-subtext">{{ formatDate(msg.date_envoi) }}</small>
                            </td>
                            <td class="td-content">
                                <span class="quote-text">"{{ msg.contenu }}"</span>
                            </td>
                            <td>
                                <span class="status-logistique" style="background: #ffe5e5; color: #d32f2f;">Signalé {{ msg.nb_signalements }}x</span>
                                <div style="margin-top: 5px;">
                                    <small class="table-subtext" v-for="(detail, i) in msg.details" :key="i">
                                        <b>{{ detail.reporter_name }}</b> le {{ formatDate(detail.date_signalement) }} : {{ detail.motif }}
                                    </small>
                                </div>
                            </td>
                            <td class="actions-cell">
                                <button class="btn-remove" style="width: 100%;" @click="deleteMessage(msg.message_id)">Supprimer</button>
                                <button class="btn-ban" style="width: 100%;" @click="toggleBan(msg.author_id, msg.author_name, true)">Bannir</button>
                                <button class="btn-ignore" style="width: 100%;" @click="ignoreMessage(msg.message_id)">Ignorer</button>                            </td>
                        </tr>
                    </tbody>
                </table>
                <div v-else class="state-card">Aucun message n'a été signalé.</div>
            </div>

            <div v-else-if="activeTab === 'topics'">
                <table v-if="allTopics.length > 0" class="data-table">
                    <thead>
                        <tr>
                            <th>SUJET</th>
                            <th>CRÉATEUR</th>
                            <th>NB MESSAGES</th>
                            <th>DATE</th>
                            <th>ACTIONS</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="topic in allTopics" :key="topic.id">
                            <td>
                                <strong>{{ topic.titre }}</strong><br />
                                <small class="table-subtext">{{ truncateText(topic.sujet, 50) }}</small>
                            </td>
                            <td class="user-link" @click="viewProfile(topic.author_id)">
                                <strong>{{ topic.author_name }}</strong>
                            </td>
                            <td><span class="badge badge--green">{{ topic.msg_count }} messages</span></td>
                            <td>{{ formatDate(topic.date) }}</td>
                            <td class="actions-cell">
                                <button class="btn-remove" @click="deleteTopic(topic.id, topic.titre)">Supprimer</button>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <div v-else class="state-card">Aucune discussion ouverte.</div>
            </div>

            <div v-else-if="activeTab === 'messages'">
                <table v-if="recentMessages.length > 0" class="data-table">
                    <thead>
                        <tr>
                            <th>DISCUSSION</th> <th>AUTEUR</th>
                            <th>MESSAGE</th>
                            <th>DATE</th>
                            <th>ACTIONS</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="msg in recentMessages" :key="msg.id">
                            <td>
                                <strong>{{ msg.topic_title }}</strong>
                            </td>
                            
                            <td class="user-link" @click="viewProfile(msg.user_id)">
                                <strong>{{ msg.author }}</strong><br />
                                <span class="status-logistique">{{ msg.role }}</span>
                            </td>
                            <td class="td-content">
                                <span class="quote-text">"{{ msg.content }}"</span>
                            </td>
                            <td>{{ formatDate(msg.postedAt) }}</td>
                            <td class="actions-cell">
                                <button class="btn-view" @click="deleteMessage(msg.id)" style="color: #d32f2f; border-color: #ffcccc;">Retirer</button>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <div v-else class="state-card">Aucun message récent.</div>
            </div>

            <div v-else-if="activeTab === 'bannis'">
                <table v-if="bannedUsers.length > 0" class="data-table">
                    <thead>
                        <tr>
                            <th>UTILISATEUR</th>
                            <th>RÔLE</th>
                            <th>STATUT</th>
                            <th>ACTIONS</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="user in bannedUsers" :key="user.id">
                            <td class="user-link" @click="viewProfile(user.id)">
                                <strong>{{ user.full_name }}</strong>
                            </td>
                            <td><span class="status-logistique">{{ user.role }}</span></td>
                            <td><span class="badge badge--red">Banni du forum</span></td>
                            <td class="actions-cell">
                                <button class="btn-view" style="color: #2d7a4f; border-color: #2d7a4f;" @click="toggleBan(user.id, user.full_name, false)">
                                    Rétablir l'accès
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <div v-else class="state-card">Aucun membre banni actuellement.</div>
            </div>

        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

const API_URL = "http://localhost:8081";
const token = sessionStorage.getItem("userToken");

const activeTab = ref('signalements');
const reportedMessages = ref([]);
const recentMessages = ref([]);
const bannedUsers = ref([]);
const allTopics = ref([]);
const loading = ref(true);
const isRefreshing = ref(false);

const headers = {
    "Content-Type": "application/json",
    "Authorization": token || ""
};

const viewProfile = (id) => {
    if (id) {
        router.push(`/user/${id}`);
    } else {
        console.error("ID utilisateur manquant");
    }
};

const fetchData = async () => {
    try {
        const [repRes, recRes, banRes, topRes] = await Promise.all([
            fetch(`${API_URL}/api/moderation/forums/signales`, { headers }),
            fetch(`${API_URL}/forums/messages/recent`, { headers }),
            fetch(`${API_URL}/api/moderation/users/banned`, { headers }),
            fetch(`${API_URL}/api/moderation/topics`, { headers })
        ]);

        if (repRes.ok) reportedMessages.value = await repRes.json();
        if (recRes.ok) recentMessages.value = await recRes.json();
        if (banRes.ok) bannedUsers.value = await banRes.json();
        if (topRes.ok) allTopics.value = await topRes.json();
    } catch (e) {
        console.error("Erreur de fetch", e);
    }
};

const refreshData = async () => {
    isRefreshing.value = true;
    await fetchData();
    setTimeout(() => isRefreshing.value = false, 500);
};

onMounted(async () => {
    await fetchData();
    loading.value = false;
});

const deleteMessage = async (msgId) => {
    if (!confirm("Voulez-vous vraiment supprimer ce message ? L'auteur recevra une notification.")) return;
    try {
        const res = await fetch(`${API_URL}/forums/message/${msgId}`, { method: "DELETE", headers });
        if (res.ok) {
            reportedMessages.value = reportedMessages.value.filter(m => m.message_id !== msgId);
            recentMessages.value = recentMessages.value.filter(m => m.id !== msgId);
        }
    } catch (e) { console.error(e); }
};

const ignoreMessage = async (msgId) => {
    if (!confirm("Voulez-vous vraiment ignorer ce signalement ? Le message sera conservé sur le forum.")) return;
    
    try {
        const res = await fetch(`${API_URL}/forums/signalement/${msgId}`, { 
            method: "DELETE", 
            headers 
        });
        
        if (res.ok) {
            reportedMessages.value = reportedMessages.value.filter(m => m.message_id !== msgId);
        } else {
            const errorMsg = await res.text();
            alert("Erreur côté serveur : " + errorMsg);
        }
    } catch (e) { 
        console.error("Erreur réseau:", e); 
    }
};

const deleteTopic = async (topicId, topicTitre) => {
    if (!confirm(`ATTENTION : Supprimer la discussion "${topicTitre}" et TOUS ses messages ?`)) return;
    try {
        const res = await fetch(`${API_URL}/forums/topic/${topicId}`, { method: "DELETE", headers });
        if (res.ok) {
            allTopics.value = allTopics.value.filter(t => t.id !== topicId);
            await fetchData();
        }
    } catch (e) { console.error(e); }
};

const toggleBan = async (userId, userName, isBanning) => {
    const actionText = isBanning ? `bannir "${userName}" du forum` : `rétablir l'accès de "${userName}"`;
    if (!confirm(`Voulez-vous vraiment ${actionText} ?`)) return;

    try {
        const res = await fetch(`${API_URL}/api/moderation/user/${userId}/ban-forum`, {
            method: "PUT",
            headers,
            body: JSON.stringify({ ban: isBanning })
        });
        if (res.ok) {
            await fetchData(); 
        }
    } catch (e) { console.error(e); }
};

const formatDate = (val) => {
    if (!val) return "NULL";
    
    const cleanDate = val.replace(/Z$/, ''); 
    
    const date = new Date(cleanDate);
    return isNaN(date.getTime())
        ? "NULL"
        : new Intl.DateTimeFormat("fr-FR", {
              day: "2-digit",
              month: "short",
              year: "numeric",
              hour: "2-digit",
              minute: "2-digit"
          }).format(date).replace(":", "h");
};

const truncateText = (text, length) => {
    if (!text) return "";
    return text.length > length ? text.substring(0, length) + '...' : text;
};
</script>

<style scoped>
.public-dashboard {
    min-height: 100vh;
    padding: 20px;
    background: var(--bg-light, #f7f9f7);
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 2rem;
}

.sidebar__category2 {
    font-size: 0.65rem;
    color: #8fa396;
    letter-spacing: 1px;
    margin: 0 0 0.5rem 0;
    text-transform: uppercase;
}

.hero-title1 {
    font-size: 2rem;
    font-weight: 800;
    margin: 1.5rem 0 0.5rem;
    color: #1a1a1a;
}

.classic-text {
    color: var(--text-grey, #666);
    margin: 0;
}

.btn-main-action {
    display: inline-flex;
    align-items: center;
    background: #2d7a4f;
    color: white;
    padding: 10px 20px;
    border-radius: 10px;
    text-decoration: none;
    font-weight: bold;
    border: none;
    cursor: pointer;
}
.btn-main-action:disabled {
    opacity: 0.7;
    cursor: not-allowed;
}

.section-container {
    background: white;
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid #eee;
    margin-bottom: 2rem;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
}

.state-card {
    border: 1px dashed #cfe0d4;
    border-radius: 14px;
    padding: 26px;
    color: var(--text-grey, #666);
    background: #fbfdfb;
    text-align: center;
}

.stats-grid {
    display: grid;
    gap: 1.5rem;
    margin-bottom: 2rem;
}

.dashboard-stats-mod {
    grid-template-columns: repeat(4, 1fr);
}

@media (max-width: 920px) {
    .dashboard-stats-mod {
        grid-template-columns: repeat(2, 1fr);
    }
}
@media (max-width: 500px) {
    .dashboard-stats-mod {
        grid-template-columns: 1fr;
    }
}

.card {
    background: white;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #eee;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
}

.card-danger-alert {
    border-color: #ffcccc;
    background: #fffafa;
}

.card-num {
    font-size: 2.5rem;
    font-weight: bold;
    color: #1a1a1a;
    line-height: 1;
    margin-bottom: 0.5rem;
}

.text-dm {
    margin-bottom: 0.8rem;
}

.badge {
    padding: 4px 10px;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: bold;
    display: inline-block;
}

.badge--green { background: #eaf4ed; color: #2d7a4f; }
.badge--orange { background: #fff3cd; color: #856404; }
.badge--red { background: #ffe5e5; color: #d32f2f; }

.tabs-actions {
    display: flex;
    gap: 10px;
}

.btn-tab {
    background: #f0f4f1;
    border: 1px solid #ddd;
    color: #666;
    padding: 8px 16px;
    border-radius: 8px;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.2s;
}

.btn-tab.active {
    background: #2d7a4f;
    color: white;
    border-color: #2d7a4f;
}

.btn-ban {
    background: #fff3cd;
    border: 1px solid #ffeeba;
    color: #856404;
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.2s;
}
.btn-ban:hover { background: #ffe8a1; color: #664d03; }

.btn-remove {
    background: #ffe5e5;
    border: 1px solid #ffcccc;
    color: #d32f2f;
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.2s;
}
.btn-remove:hover { background: #ffcccc; color: #b71c1c; }

.btn-view {
    background: transparent;
    border: 1px solid #ddd;
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
}
.btn-view:hover { background: #f0f0f0; }

.status-logistique {
    background: #f0f0f0;
    color: #555;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: bold;
    text-transform: uppercase;
}

.data-table {
    width: 100%;
    border-collapse: collapse;
}

.data-table th {
    text-align: left;
    padding: 12px;
    font-size: 0.85rem;
    color: #666;
    border-bottom: 2px solid #eee;
}

.data-table td {
    padding: 16px 12px;
    border-bottom: 1px solid #eee;
    vertical-align: middle;
}

.td-content {
    max-width: 400px;
}

.quote-text {
    font-style: italic;
    color: #444;
    background: #f9f9f9;
    padding: 6px 10px;
    border-radius: 6px;
    border-left: 3px solid #ccc;
    display: block;
}

.table-subtext {
    display: block;
    margin-top: 4px;
    color: var(--text-grey, #666);
    font-size: 0.78rem;
}

.actions-cell {
    text-align: right;
    min-width: 120px;
}

.user-link {
    cursor: pointer;
    transition: all 0.2s ease;
}

.user-link strong {
    transition: color 0.2s ease;
}

.user-link:hover strong {
    color: #2d7a4f; 
}

.btn-ignore {
    margin-top: 0;
    background: #f5f5f5;
    border: 1px solid #ddd;
    color: #666;
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.2s;
}

.btn-ignore:hover { 
    background: #e0e0e0; 
    color: #333; 
}
</style>