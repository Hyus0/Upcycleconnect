<template>
    <div class="info-view">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    ACCUEIL > PARAMÈTRES > NOTIFICATIONS
                </p>
                <h1 class="hero-title1">MES NOTIFICATIONS</h1>
                <p class="classic-text">
                    Consultez vos alertes casiers, likes et mises à jour de
                    projets.
                </p>

                <div
                    v-if="errors.length > 0"
                    style="
                        background-color: #fee2e2;
                        border: 1px solid #ef4444;
                        color: #b91c1c;
                        padding: 10px;
                        border-radius: 8px;
                        margin-bottom: 15px;
                    "
                >
                    <ul style="margin: 0; padding-left: 20px">
                        <li v-for="(err, index) in errors" :key="index">
                            {{ err }}
                        </li>
                    </ul>
                </div>

                <div
                    v-if="successMsg"
                    class="success-box"
                    style="
                        background-color: #e2fee3;
                        border: 1px solid #44ef44;
                        color: #158f3c;
                        padding: 10px;
                        border-radius: 8px;
                        margin-bottom: 15px;
                    "
                >
                    {{ successMsg }}
                </div>
            </div>

            <button
                v-if="unreadCount > 0"
                class="btn-main-action"
                @click="marquerToutCommeLu"
            >
                ✓ Tout marquer comme lu ({{ unreadCount }})
            </button>
        </header>

        <div class="info-layout">
            <div class="info-main-col">
                <section class="info-section">
                    <h2 class="section-title">Boîte de réception</h2>

                    <div v-if="loading" class="empty-state">
                        Chargement de vos notifications...
                    </div>

                    <div
                        v-else-if="notifications.length === 0"
                        class="empty-state"
                    >
                        Vous n'avez aucune notification pour le moment. 🍃
                    </div>

                    <div v-else class="notifications-list">
                        <div
                            v-for="notif in notifications"
                            :key="notif.id"
                            class="notif-card"
                            :class="{ 'is-unread': !notif.lu }"
                            @click="marquerCommeLu(notif)"
                        >
                            <div class="notif-icon">
                                <component
                                    :is="getIconComponent(notif.type)"
                                    :color="getIconColor(notif.type)"
                                    :size="24"
                                />
                            </div>

                            <div class="notif-content">
                                <div class="notif-header">
                                    <h3 class="notif-title">
                                        {{ notif.titre }}
                                    </h3>
                                    <span class="notif-date">{{
                                        formatDate(notif.date_envoi)
                                    }}</span>
                                </div>
                                <p class="notif-message">{{ notif.message }}</p>
                            </div>

                            <div class="notif-status" v-if="!notif.lu">
                                <span class="unread-dot"></span>
                            </div>
                        </div>
                    </div>
                </section>
            </div>

            <aside class="info-side-col">
                <section class="info-section status-card">
                    <h2 class="section-title">Résumé des alertes</h2>
                    <div class="status-item">
                        <span class="status-label">Total reçues</span>
                        <span class="badge-role">{{
                            notifications.length
                        }}</span>
                    </div>
                    <div class="status-item">
                        <span class="status-label">Non lues</span>
                        <span
                            class="status-active"
                            :style="{
                                color: unreadCount > 0 ? '#ef4444' : '#1e7e34',
                            }"
                        >
                            {{ unreadCount }}
                        </span>
                    </div>
                    <div class="divider"></div>
                    <div class="registration-info">
                        <p>
                            Gardez un œil sur cette section pour ne manquer
                            aucune interaction avec la communauté !
                        </p>
                    </div>
                </section>
            </aside>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";

import {
    TriangleAlert,
    Mail,
    CircleAlert,
    Cuboid,
    ThumbsUp,
    MessageSquare,
    Rss,
    Bell,
} from "lucide-vue-next";

const API_URL = "/go";
const notifications = ref([]);
const loading = ref(true);
const errors = ref([]);
const successMsg = ref("");

const unreadCount = computed(() => {
    return notifications.value.filter((n) => !n.lu).length;
});

const getIconComponent = (type) => {
    if (!type) return Bell;
    switch (type.toLowerCase()) {
        case "alerte":
            return TriangleAlert;
        case "message":
            return Mail;
        case "rappel":
            return CircleAlert;
        case "casier":
            return Cuboid;
        case "like":
            return ThumbsUp;
        case "avis":
            return MessageSquare;
        case "follow":
            return Rss;
        default:
            return Bell;
    }
};

const getIconColor = (type) => {
    if (!type) return "#6b7280";
    switch (type.toLowerCase()) {
        case "alerte":
            return "#ef4444";
        case "message":
            return "#3b82f6";
        case "rappel":
            return "#f59e0b";
        case "casier":
            return "#8b5cf6";
        case "like":
            return "#ec4899";
        case "avis":
            return "#10b981";
        case "follow":
            return "#06b6d4";
        default:
            return "#6b7280";
    }
};

const formatDate = (dateString) => {
    if (!dateString) return "...";
    const date = new Date(dateString);
    return date.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
        hour: "2-digit",
        minute: "2-digit",
    });
};

const fetchNotifications = async () => {
    const userId = sessionStorage.getItem("userId");
    const token = sessionStorage.getItem("userToken");

    if (!userId || !token) {
        errors.value = ["Session expirée. Veuillez vous reconnecter."];
        loading.value = false;
        return;
    }

    loading.value = true;
    try {
        const res = await fetch(`${API_URL}/users/${userId}/notifications`, {
            method: "GET",
            headers: {
                Authorization: token,
                "Content-Type": "application/json",
            },
        });

        if (res.ok) {
            notifications.value = await res.json();
        } else {
            errors.value = ["Impossible de charger les notifications."];
        }
    } catch (error) {
        errors.value = ["Erreur réseau lors de la récupération."];
    } finally {
        loading.value = false;
    }
};

const marquerCommeLu = async (notif) => {
    if (notif.lu) return;

    const token = sessionStorage.getItem("userToken");
    const userId =
        sessionStorage.getItem("userId") || sessionStorage.getItem("id");

    try {
        const res = await fetch(
            `/go/notifications/${notif.id}/read`,
            {
                method: "POST",
                headers: {
                    Authorization: token,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ user_id: Number(userId) }),
            },
        );

        if (res.ok) {
            notif.lu = true;
        } else {
            console.error("Erreur serveur :", await res.text());
        }
    } catch (err) {
        console.error("Erreur réseau :", err);
    }
};

const marquerToutCommeLu = async () => {
    const unreadNotifs = notifications.value.filter((n) => !n.lu);
    if (unreadNotifs.length === 0) return;

    for (let notif of unreadNotifs) {
        await marquerCommeLu(notif);
    }
    successMsg.value = "Toutes les notifications ont été marquées comme lues.";
    setTimeout(() => {
        successMsg.value = "";
    }, 3000);
};

onMounted(() => {
    fetchNotifications();
});
</script>

<style scoped>
.info-view {
    padding: 10px;
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 2rem;
}

.header-left {
    display: flex;
    flex-direction: column;
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
    color: #666;
    margin: 0 0 15px 0;
}

.btn-main-action {
    background-color: #2d7a4f;
    color: white;
    border: none;
    padding: 12px 20px;
    border-radius: 10px;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s;
    white-space: nowrap;
}

.btn-main-action:hover {
    background-color: #246343;
}

.info-layout {
    display: grid;
    grid-template-columns: 1fr 320px;
    gap: 25px;
    margin-top: 20px;
}

.info-section {
    background: white;
    border-radius: 16px;
    padding: 24px;
    margin-bottom: 25px;
    border: 1px solid #e8ebe9;
}

.section-title {
    font-size: 1.1rem;
    font-family: "Syne", sans-serif;
    font-weight: 700;
    margin-bottom: 20px;
    color: #1a1f1c;
}

.empty-state {
    text-align: center;
    padding: 3rem 1rem;
    color: #8fa396;
    font-style: italic;
    font-size: 1rem;
}

.notifications-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.notif-card {
    display: flex;
    align-items: flex-start;
    gap: 15px;
    padding: 16px;
    border-radius: 12px;
    border: 1px solid #f0f4f1;
    background: #fcfdfc;
    cursor: pointer;
    transition: all 0.2s;
}

.notif-card:hover {
    border-color: #dcdfdc;
    background: #ffffff;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.is-unread {
    background: #f0fdf4;
    border: 1px solid #dcfce7;
}
.is-unread:hover {
    background: #eaf4ed;
}

.notif-icon {
    background: white;
    width: 45px;
    height: 45px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    flex-shrink: 0;
    border: 1px solid #e8ebe9;
}
.is-unread .notif-icon {
    border-color: #bbf7d0;
}

.notif-content {
    flex-grow: 1;
}

.notif-header {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    margin-bottom: 4px;
}

.notif-title {
    margin: 0;
    font-size: 1rem;
    font-weight: 600;
    color: #1a1f1c;
}
.is-unread .notif-title {
    font-weight: 700;
    color: #166534;
}

.notif-date {
    font-size: 0.75rem;
    color: #a0ada7;
    white-space: nowrap;
    margin-left: 10px;
}

.notif-message {
    margin: 0;
    font-size: 0.9rem;
    color: #5a6660;
    line-height: 1.5;
}

.notif-status {
    display: flex;
    align-items: center;
    justify-content: center;
    padding-top: 5px;
}

.unread-dot {
    width: 10px;
    height: 10px;
    background-color: #22c55e;
    border-radius: 50%;
}

.status-item {
    display: flex;
    justify-content: space-between;
    margin-bottom: 15px;
    font-size: 0.9rem;
}

.status-label {
    color: #5a6660;
}

.badge-role {
    background-color: #f0f4f1;
    color: #2d7a4f;
    padding: 4px 10px;
    border-radius: 8px;
    font-weight: bold;
    font-size: 0.8rem;
}

.status-active {
    font-weight: bold;
}

.divider {
    height: 1px;
    background: #e8ebe9;
    margin: 20px 0;
}

.registration-info p {
    font-size: 0.85rem;
    color: #5a6660;
    margin: 0;
    line-height: 1.5;
}

@media (max-width: 1024px) {
    .info-layout {
        grid-template-columns: 1fr;
    }
    .info-side-col {
        order: -1;
    }
}

</style>
