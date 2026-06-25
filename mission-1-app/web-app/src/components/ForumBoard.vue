<template>
    <section class="forum-page">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > COMMUNAUTE > FORUMS</p>
                <h1 class="hero-title1">FORUMS COMMUNAUTAIRES</h1>
                <p class="classic-text">
                    Retrouvez les discussions utiles, posez vos questions et
                    partagez vos retours d'experience.
                </p>
            </div>
            <button
                v-if="!isCheckingBan && !isBanned"
                class="btn-main-action"
                type="button"
                @click="openCreateTopic"
            >
                + Nouvelle discussion
            </button>
        </header>

        <div v-if="!isAuthenticated" class="forum-guest-banner">
            <div>
                <span class="forum-chip">Lecture ouverte</span>
                <h2>Vous pouvez consulter tous les forums sans compte.</h2>
                <p>
                    Pour publier une discussion ou repondre, connectez-vous ou
                    creez votre compte UpcycleConnect.
                </p>
            </div>
            <div class="forum-guest-banner__actions">
                <RouterLink
                    class="btn-secondary forum-guest-button"
                    to="/connexion"
                    >Connexion</RouterLink
                >
                <RouterLink
                    class="btn-main-action forum-guest-button"
                    to="/inscription"
                    >Creer un compte</RouterLink
                >
            </div>
        </div>

        <div v-if="isAuthenticated && isBanned" class="forum-guest-banner" style="background: #fceaea; border-color: #f5c6c6;">
            <div>
                <span class="forum-chip" style="background: #e74c3c; color: white;">Accès Restreint</span>
                <h2 style="color: #c0392b;">Vous avez été restreint de l'espace communautaire.</h2>
                <p style="color: #e74c3c;">
                    Suite à une modération, vous pouvez toujours lire les discussions, mais vous n'avez plus l'autorisation de publier des messages ou de créer de nouveaux sujets.
                </p>
            </div>
        </div>

        <div class="forum-overview">
            <article class="forum-stat forum-stat--primary">
                <span class="forum-stat__label">Forums actifs</span>
                <strong>{{ forums.length }}</strong>
                <small>Salons communautaires disponibles</small>
            </article>
            <article class="forum-stat">
                <span class="forum-stat__label">Discussions</span>
                <strong>{{ totalTopics }}</strong>
                <small>Fils ouverts sur votre espace</small>
            </article>
            <article class="forum-stat">
                <span class="forum-stat__label">Messages</span>
                <strong>{{ totalMessages }}</strong>
                <small>Contributions publiees</small>
            </article>
        </div>

        <div class="forum-layout">
            <aside class="section-container forum-sidebar">
                <div class="forum-sidebar__head">
                    <div>
                        <span class="forum-chip">Salons</span>
                        <h2>Explorer</h2>
                    </div>
                    <span class="forum-count">{{ forums.length }}</span>
                </div>

                <div class="forum-sidebar__search">
                    <input
                        v-model.trim="search"
                        type="text"
                        placeholder="Rechercher un sujet ou un salon"
                        class="search-input forum-search"
                    />
                </div>

                <div class="forum-list">
                    <button
                        v-for="forum in filteredForums"
                        :key="forum.id"
                        type="button"
                        class="forum-card"
                        :class="{ 'is-active': forum.id === selectedForumId }"
                        @click="selectForum(forum.id)"
                    >
                        <div class="forum-card__top">
                            <strong>{{ forum.name }}</strong>
                            <span>{{
                                forum.topics ? forum.topics.length : 0
                            }}</span>
                        </div>
                        <p>{{ forum.description }}</p>
                        <small>{{ latestTopicLabel(forum) }}</small>
                    </button>
                </div>
            </aside>

            <div class="forum-content">
                <article
                    class="section-container forum-main"
                    v-if="selectedForum"
                >
                    <div class="forum-main__head">
                        <div>
                            <span class="forum-chip">{{
                                selectedForum.name
                            }}</span>
                            <h2>{{ selectedForum.description }}</h2>
                        </div>
                        <button
                            v-if="!isCheckingBan && !isBanned"
                            class="btn-secondary"
                            type="button"
                            @click="openCreateTopic"
                        >
                            Ecrire un sujet
                        </button>
                    </div>

                    <div v-if="showComposer && !isBanned" class="forum-composer">
                        <div class="forum-composer__head">
                            <strong>{{
                                selectedTopic
                                    ? "Repondre a la discussion"
                                    : "Nouvelle discussion"
                            }}</strong>
                            <button
                                class="forum-link"
                                type="button"
                                @click="closeComposer"
                            >
                                Annuler
                            </button>
                        </div>
                        <input
                            v-if="!selectedTopic"
                            v-model.trim="draftTopicTitle"
                            type="text"
                            placeholder="Titre de la discussion"
                            class="forum-input"
                        />
                        <textarea
                            v-model.trim="draftMessage"
                            placeholder="Ecrivez votre message pour la communaute..."
                            class="forum-textarea"
                            rows="5"
                        ></textarea>
                        <div class="forum-composer__actions">
                            <span class="forum-helper"
                                >Publie sous le nom {{ userName }}</span
                            >
                            <button
                                class="btn-main-action"
                                type="button"
                                @click="submitPost"
                            >
                                {{
                                    selectedTopic
                                        ? "Envoyer la reponse"
                                        : "Publier la discussion"
                                }}
                            </button>
                        </div>
                    </div>

                    <div class="forum-topics">
                        <article
                            v-for="topic in filteredTopics"
                            :key="topic.id"
                            class="forum-topic"
                            :class="{ 'is-open': topic.id === selectedTopicId }"
                        >
                            <button
                                type="button"
                                class="forum-topic__summary"
                                @click="toggleTopic(topic.id)"
                            >
                                <div>
                                    <span
                                        class="forum-chip forum-chip--muted"
                                        >{{ topic.tag || "Discussion" }}</span
                                    >
                                    <h3>{{ topic.title }}</h3>
                                    <p>{{ topic.preview }}</p>
                                </div>
                                <div class="forum-topic__meta">
                                    <strong
                                        >{{
                                            topic.messages
                                                ? topic.messages.length
                                                : 0
                                        }}
                                        msg</strong
                                    >
                                    <small>{{ formatDate(topic.lastActivity) }}</small>
                                </div>
                            </button>

                            <div
                                v-if="topic.id === selectedTopicId"
                                class="forum-thread"
                            >
                                <div
                                    v-for="message in topic.messages"
                                    :key="message.id"
                                    class="forum-message"
                                    :class="{
                                        'is-me': message.author === userName,
                                    }"
                                >
                                    <div
                                        class="forum-message__avatar clickable-user"
                                        @click="
                                            goToUserProfile(
                                                message.user_id ||
                                                    message.id_auteur,
                                            )
                                        "
                                    >
                                        <img
                                            v-if="
                                                message.image_profil &&
                                                message.image_profil.trim() !==
                                                    ''
                                            "
                                            :src="message.image_profil"
                                            alt="Avatar"
                                            class="forum-message__avatar-img"
                                        />
                                        <span v-else>
                                            {{ initialsFor(message.author) }}
                                        </span>
                                    </div>
                                    <div class="forum-message__body">
                                        <div class="forum-message__meta">
                                            <div class="meta-left">
                                                <strong
                                                    class="clickable-user"
                                                    @click="
                                                        goToUserProfile(
                                                            message.user_id ||
                                                                message.id_auteur,
                                                        )
                                                    "
                                                >
                                                    {{ message.author }}
                                                </strong>
                                                <span>{{ message.role }}</span>
                                                <small>{{ formatDate(message.postedAt) }}</small>
                                            </div>
                                            <button
                                                v-if="message.author !== userName && !isBanned"
                                                type="button"
                                                class="btn-report"
                                                title="Signaler ce message"
                                                @click="openReportModal(message.id)"
                                            >
                                                <TriangleAlert :size="16" />
                                            </button>
                                        </div>
                                        <p>{{ message.content }}</p>
                                    </div>
                                </div>

                                <div class="forum-thread__footer">
                                    <button
                                        v-if="!isCheckingBan && !isBanned"
                                        class="btn-view"
                                        type="button"
                                        @click="replyToTopic(topic.id)"
                                    >
                                        Repondre
                                    </button>
                                </div>
                            </div>
                        </article>
                    </div>

                    <div v-if="filteredTopics.length === 0" class="forum-empty">
                        <p>
                            Aucune discussion ne correspond a votre recherche.
                        </p>
                        <button
                            v-if="!isCheckingBan && !isBanned"
                            class="btn-secondary"
                            type="button"
                            @click="openCreateTopic"
                        >
                            Creer un sujet
                        </button>
                    </div>
                </article>
            </div>
        </div>

        <div
            v-if="showReportModal"
            class="modal-overlay"
            @click.self="closeReportModal"
        >
            <div class="modal-content">
                <h3>Signaler un message</h3>
                <p>
                    Veuillez indiquer le motif de votre signalement pour aider
                    nos modérateurs :
                </p>
                <textarea
                    v-model="reportMotif"
                    rows="4"
                    placeholder="Ex: Spam, insulte, hors-sujet..."
                    class="forum-textarea"
                ></textarea>
                <div class="modal-actions">
                    <button
                        class="btn-secondary"
                        type="button"
                        @click="closeReportModal"
                    >
                        Annuler
                    </button>
                    <button
                        class="btn-action-danger"
                        type="button"
                        @click="submitReport"
                    >
                        Envoyer
                    </button>
                </div>
            </div>
        </div>
    </section>
</template>

<script setup>
import { computed, ref, onMounted } from "vue";
import { RouterLink, useRouter } from "vue-router";
import { TriangleAlert } from "lucide-vue-next"; 

const router = useRouter();
const API_URL = "/go";

const userId = sessionStorage.getItem("userId") || "guest";
const userPrenom = sessionStorage.getItem("userPrenom") || "";
const userNom = sessionStorage.getItem("userNom") || "";
const userName = `${userPrenom} ${userNom}`.trim() || "Utilisateur";

const isAuthenticated = computed(() =>
    Boolean(
        sessionStorage.getItem("userToken") && sessionStorage.getItem("userId"),
    ),
);

const forums = ref([]);
const search = ref("");
const selectedForumId = ref("");
const selectedTopicId = ref("");
const showComposer = ref(false);
const draftTopicTitle = ref("");
const draftMessage = ref("");

const showReportModal = ref(false);
const reportMessageId = ref(null);
const reportMotif = ref("");

// État de ban de l'utilisateur
const isBanned = ref(false);
const isCheckingBan = ref(true); // Ajout pour empêcher le clignotement des boutons

const checkBanStatus = async () => {
    if (!isAuthenticated.value) {
        isCheckingBan.value = false;
        return;
    }
    
    try {
        const currentUserId = parseInt(sessionStorage.getItem("userId"), 10);
        
        // ASTUCE: On utilise directement l'API des utilisateurs bannis
        const res = await fetch(`${API_URL}/api/moderation/users/banned`);
        if (res.ok) {
            const bannedUsers = await res.json() || [];
            
            // On vérifie si l'ID de l'utilisateur actuel est dans la liste des bannis
            isBanned.value = bannedUsers.some(user => user.id === currentUserId);
        }
    } catch (error) {
        console.error("Erreur lors de la vérification du statut:", error);
    } finally {
        // La vérification est terminée, on peut afficher ou cacher les boutons
        isCheckingBan.value = false;
    }
};

const fetchForums = async () => {
    try {
        const res = await fetch(`${API_URL}/forums`);
        if (res.ok) {
            forums.value = (await res.json()) || [];
        }
    } catch (error) {
        console.error(error);
    } finally {
        if (forums.value.length === 0) {
            forums.value = [
                {
                    id: "general",
                    name: "Discussions Générales",
                    description: "Toutes les discussions de la communauté",
                    topics: [],
                },
            ];
        }

        if (!selectedForumId.value) {
            selectedForumId.value = forums.value[0].id;
            if (forums.value[0].topics && forums.value[0].topics.length > 0) {
                selectedTopicId.value = forums.value[0].topics[0].id;
            }
        }
    }
};

onMounted(() => {
    fetchForums();
    checkBanStatus(); // Appel lors du chargement de la page
});

function formatDate(dateString) {
    if (!dateString) return "";
    const date = new Date(dateString);
    
    if (isNaN(date.getTime())) return dateString; 
    
    return new Intl.DateTimeFormat("fr-FR", {
        day: "2-digit",
        month: "short",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit"
    }).format(date).replace(":", "h");
}

const selectedForum = computed(
    () =>
        forums.value.find((forum) => forum.id === selectedForumId.value) ||
        null,
);
const selectedTopic = computed(
    () =>
        selectedForum.value?.topics?.find(
            (topic) => topic.id === selectedTopicId.value,
        ) || null,
);

const filteredForums = computed(() => {
    const term = search.value.toLowerCase();
    if (!term) return forums.value;

    return forums.value.filter((forum) => {
        const forumMatch =
            forum.name.toLowerCase().includes(term) ||
            forum.description.toLowerCase().includes(term);
        const topics = forum.topics || [];
        const topicMatch = topics.some(
            (topic) =>
                topic.title.toLowerCase().includes(term) ||
                topic.preview.toLowerCase().includes(term),
        );
        return forumMatch || topicMatch;
    });
});

const filteredTopics = computed(() => {
    if (!selectedForum.value || !selectedForum.value.topics) return [];
    const term = search.value.toLowerCase();
    if (!term) return selectedForum.value.topics;

    return selectedForum.value.topics.filter(
        (topic) =>
            topic.title.toLowerCase().includes(term) ||
            topic.preview.toLowerCase().includes(term) ||
            (topic.messages || []).some((message) =>
                message.content.toLowerCase().includes(term),
            ),
    );
});

const totalTopics = computed(() =>
    forums.value.reduce(
        (count, forum) => count + (forum.topics ? forum.topics.length : 0),
        0,
    ),
);
const totalMessages = computed(() =>
    forums.value.reduce(
        (count, forum) =>
            count +
            (forum.topics
                ? forum.topics.reduce(
                      (topicCount, topic) =>
                          topicCount +
                          (topic.messages ? topic.messages.length : 0),
                      0,
                  )
                : 0),
        0,
    ),
);

function latestTopicLabel(forum) {
    if (!forum.topics || !forum.topics.length) return "Aucune discussion";
    return forum.topics[0].title;
}

function goToUserProfile(targetId) {
    if (targetId) {
        router.push(`/user/${targetId}`);
    } else {
        alert(
            "Impossible d'ouvrir le profil : l'ID de l'utilisateur n'est pas fourni par l'API.",
        );
    }
}

function initialsFor(name) {
    if (!name) return "?";
    return name
        .split(" ")
        .filter(Boolean)
        .slice(0, 2)
        .map((part) => part[0]?.toUpperCase() || "")
        .join("");
}

// Fonction de sécurité JS qui empêche l'action même si l'interface est forcée
function requireAuth() {
    if (isAuthenticated.value) {
        if (isBanned.value) {
            alert("Accès refusé : Vous êtes banni du forum et ne pouvez plus participer.");
            return false;
        }
        return true;
    }
    router.push({ path: "/connexion", query: { redirect: "/forums" } });
    return false;
}

function selectForum(forumId) {
    selectedForumId.value = forumId;
    const forum = forums.value.find((f) => f.id === forumId);
    selectedTopicId.value =
        forum && forum.topics && forum.topics.length > 0
            ? forum.topics[0].id
            : "";
    closeComposer();
}

function toggleTopic(topicId) {
    selectedTopicId.value = selectedTopicId.value === topicId ? "" : topicId;
    showComposer.value = false;
}

function openCreateTopic() {
    if (!requireAuth()) return;
    selectedTopicId.value = "";
    showComposer.value = true;
}

function replyToTopic(topicId) {
    if (!requireAuth()) return;
    selectedTopicId.value = topicId;
    showComposer.value = true;
}

function closeComposer() {
    showComposer.value = false;
    draftTopicTitle.value = "";
    draftMessage.value = "";
}

function openReportModal(messageId) {
    if (!requireAuth()) return;
    reportMessageId.value = messageId;
    reportMotif.value = "";
    showReportModal.value = true;
}

function closeReportModal() {
    showReportModal.value = false;
    reportMessageId.value = null;
    reportMotif.value = "";
}

async function submitReport() {
    if (!reportMotif.value.trim()) {
        alert("Veuillez renseigner un motif pour ce signalement.");
        return;
    }

    const currentUserId = parseInt(sessionStorage.getItem("userId"), 10);

    try {
        const res = await fetch(
            `${API_URL}/forums/message/${reportMessageId.value}/signaler`,
            {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    user_id: currentUserId,
                    motif: reportMotif.value.trim(),
                }),
            },
        );

        if (res.ok) {
            alert(
                "Signalement envoyé. L'équipe de modération va examiner ce message.",
            );
            closeReportModal();
        } else {
            const errText = await res.text();
            alert("Erreur: Vous avez probablement déjà signalé ce message.");
        }
    } catch (e) {
        console.error(e);
        alert("Erreur de connexion.");
    }
}

async function submitPost() {
    if (!requireAuth() || !selectedForum.value) return;

    if (!draftMessage.value.trim()) {
        alert("Veuillez écrire un message.");
        return;
    }

    const currentUserId = parseInt(sessionStorage.getItem("userId"), 10);

    try {
        if (selectedTopic.value) {
            const res = await fetch(`${API_URL}/forums/message`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    user_id: currentUserId,         
                    forum_id: selectedTopic.value.id, 
                    contenu: draftMessage.value.trim(),
                }),
            });
            if (!res.ok) throw new Error(await res.text());
        } else {
            if (draftTopicTitle.value.trim().length < 4) {
                alert("Le titre doit faire au moins 4 caractères.");
                return;
            }

            const res = await fetch(`${API_URL}/forums/topic`, { 
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    user_id: currentUserId,            
                    salon_id: selectedForum.value.id, 
                    title: draftTopicTitle.value.trim(), 
                    sujet: draftMessage.value.trim(),
                }),
            });
            if (!res.ok) throw new Error(await res.text());
        }

        closeComposer();
        await fetchForums();
    } catch (e) {
        console.error("Erreur backend:", e);
        alert("Erreur lors de la publication : " + e.message);
    }
}
</script>

<style scoped>
.forum-page {
    display: grid;
    gap: 24px;
}
.forum-guest-banner {
    display: flex;
    justify-content: space-between;
    gap: 20px;
    align-items: center;
    padding: 22px 24px;
    border-radius: 18px;
    background: linear-gradient(135deg, #edf7f1, #f8fbf9);
    border: 1px solid #d9e8de;
}
.forum-guest-banner h2 {
    margin: 10px 0 6px;
    font-family: "Syne", sans-serif;
    font-size: 1.35rem;
}
.forum-guest-banner p {
    margin: 0;
    color: var(--text-secondary);
}
.forum-guest-banner__actions {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
}
.forum-guest-button {
    text-decoration: none;
}
.forum-overview {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 18px;
}
.forum-stat {
    padding: 22px;
    border-radius: 18px;
    background: #ffffff;
    box-shadow: var(--shadow);
    display: grid;
    gap: 6px;
}
.forum-stat--primary {
    background: linear-gradient(135deg, #2f8353, #246742);
    color: #ffffff;
}
.forum-stat--primary .forum-stat__label,
.forum-stat--primary small {
    color: rgba(255, 255, 255, 0.76);
}
.forum-stat__label {
    color: var(--text-secondary);
    font-size: 0.78rem;
    letter-spacing: 0.12em;
    text-transform: uppercase;
}
.forum-stat strong {
    font-family: "Syne", sans-serif;
    font-size: 2.25rem;
    line-height: 1;
}
.forum-stat small {
    color: var(--text-secondary);
}
.forum-layout {
    display: grid;
    grid-template-columns: 320px minmax(0, 1fr);
    gap: 22px;
    align-items: start;
}
.forum-sidebar,
.forum-main {
    margin-bottom: 0;
}
.forum-sidebar {
    position: sticky;
    top: 110px;
    display: grid;
    gap: 18px;
}
.forum-sidebar__head,
.forum-main__head,
.forum-composer__head,
.forum-composer__actions,
.forum-card__top,
.forum-topic__summary,
.forum-message__meta {
    display: flex;
    justify-content: space-between;
    gap: 14px;
}
.forum-sidebar__head h2,
.forum-main__head h2,
.forum-topic h3 {
    margin: 8px 0 0;
    font-family: "Syne", sans-serif;
}
.forum-chip {
    display: inline-flex;
    align-items: center;
    height: 30px;
    padding: 0 12px;
    border-radius: 999px;
    background: #e7f5ec;
    color: #2c7e4f;
    font-weight: 700;
    font-size: 0.78rem;
}
.forum-chip--muted {
    background: #f0f4f1;
    color: #6d7b72;
}
.forum-count {
    width: 34px;
    height: 34px;
    border-radius: 12px;
    display: grid;
    place-items: center;
    background: #f1f7f3;
    color: var(--brand-green-deep);
    font-weight: 700;
}
.forum-search {
    width: 100%;
}
.forum-list,
.forum-topics {
    display: grid;
    gap: 12px;
}
.forum-card {
    border: 1px solid #e0ebe2;
    border-radius: 16px;
    padding: 16px;
    background: linear-gradient(180deg, #fbfdfb, #f5faf6);
    text-align: left;
    display: grid;
    gap: 8px;
}
.forum-card.is-active {
    border-color: #8fc5a2;
    box-shadow: 0 12px 24px rgba(51, 132, 84, 0.12);
}
.forum-card__top strong,
.forum-topic__summary strong {
    font-size: 1rem;
}
.forum-card__top span {
    min-width: 28px;
    height: 28px;
    border-radius: 10px;
    display: grid;
    place-items: center;
    background: #ecf7f0;
    color: var(--brand-green-deep);
    font-weight: 700;
}
.forum-card p,
.forum-card small,
.forum-helper,
.forum-empty p {
    margin: 0;
    color: var(--text-secondary);
}
.forum-content {
    min-width: 0;
}
.forum-main {
    display: grid;
    gap: 18px;
}
.forum-main__head {
    align-items: flex-start;
}
.forum-main__head h2 {
    font-size: 1.55rem;
}
.forum-composer {
    display: grid;
    gap: 14px;
    padding: 18px;
    border: 1px solid #dce8de;
    border-radius: 18px;
    background: linear-gradient(180deg, #fbfdfb, #f3f8f4);
}

.forum-input, 
.forum-textarea {
  width: 100%; 
  border: 1px solid #d8e2db; 
  border-radius: 14px; 
  padding: 14px 16px; 
  background: #ffffff; 
  color: var(--text-primary);
  font-family: inherit; 
}

.forum-textarea {
    resize: vertical;
    min-height: 138px;
}
.forum-link {
    border: 0;
    background: transparent;
    color: var(--brand-green-deep);
    font-weight: 700;
}
.forum-topic {
    border: 1px solid #e5ede7;
    border-radius: 18px;
    overflow: hidden;
    background: #ffffff;
}
.forum-topic.is-open {
    border-color: #9bcbae;
    box-shadow: 0 14px 32px rgba(28, 88, 52, 0.08);
}
.forum-topic__summary {
    width: 100%;
    padding: 18px;
    border: 0;
    background: transparent;
    align-items: flex-start;
    text-align: left;
}
.forum-topic__summary h3 {
    font-size: 1.16rem;
}
.forum-topic__meta {
    min-width: 120px;
    display: grid;
    gap: 4px;
    justify-items: end;
    color: var(--text-secondary);
}
.forum-thread {
    display: grid;
    gap: 14px;
    padding: 0 18px 18px;
}
.forum-message {
    display: grid;
    grid-template-columns: 44px minmax(0, 1fr);
    gap: 14px;
    padding: 16px;
    border-radius: 16px;
    background: #f7faf7;
}
.forum-message.is-me {
    background: #edf7f0;
}
.forum-message__avatar {
    width: 44px;
    height: 44px;
    border-radius: 14px;
    display: grid;
    place-items: center;
    background: linear-gradient(180deg, #3e985f, #30794b);
    color: #ffffff;
    font-weight: 700;
    overflow: hidden;
    flex-shrink: 0;
}
.forum-message__avatar-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}
.clickable-user {
    cursor: pointer;
    transition:
        opacity 0.2s ease,
        color 0.2s ease;
}
.clickable-user:hover {
    opacity: 0.8;
    color: #2c7e4f;
    text-decoration: underline;
    text-decoration-color: #2c7e4f;
}
.forum-message__avatar.clickable-user:hover {
    text-decoration: none;
}
.forum-message__body {
    display: grid;
    gap: 8px;
}
.forum-message__meta {
    align-items: center;
    flex-wrap: wrap;
}
.forum-message__meta span,
.forum-message__meta small {
    color: var(--text-secondary);
}
.forum-message__body p {
    margin: 0;
    line-height: 1.6;
}
.forum-thread__footer {
    display: flex;
    justify-content: flex-end;
}
.forum-empty {
    padding: 18px;
    border-radius: 16px;
    background: #f8fbf9;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
}

@media (max-width: 1080px) {
    .forum-overview,
    .forum-layout {
        grid-template-columns: 1fr;
    }
    .forum-sidebar {
        position: static;
    }
}
@media (max-width: 760px) {
    .content-header,
    .forum-main__head,
    .forum-composer__head,
    .forum-composer__actions,
    .forum-empty,
    .forum-topic__summary,
    .forum-guest-banner {
        flex-direction: column;
        align-items: flex-start;
    }
    .forum-overview {
        gap: 12px;
    }
    .forum-message {
        grid-template-columns: 1fr;
    }
}

.meta-left {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
}

.btn-report {
    background: transparent;
    border: none;
    color: #a0aec0;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 4px;
    border-radius: 6px;
    transition: all 0.2s;
}

.btn-report:hover {
    color: #e74c3c;
    background: #fceaea;
}

.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.4);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(2px);
}

.modal-content {
    background: white;
    padding: 24px;
    border-radius: 16px;
    width: 90%;
    max-width: 450px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
    display: grid;
    gap: 16px;
}

.modal-content h3 {
    margin: 0;
    font-size: 1.3rem;
    color: #1a1a1a;
    font-family: "Syne", sans-serif;
}

.modal-content p {
    margin: 0;
    color: #666;
    font-size: 0.95rem;
}

.modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 8px;
}

.btn-action-danger {
    background-color: #e74c3c;
    color: white;
    border: none;
    padding: 10px 18px;
    border-radius: 10px;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s;
}

.btn-action-danger:hover {
    background-color: #c0392b;
}
</style>