<template>
  <section class="forum-page">
    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">ACCUEIL > COMMUNAUTE > FORUMS</p>
        <h1 class="hero-title1">FORUMS COMMUNAUTAIRES</h1>
        <p class="classic-text">
          Retrouvez les discussions utiles, posez vos questions et partagez vos retours d'experience.
        </p>
      </div>
      <button class="btn-main-action" type="button" @click="openCreateTopic">
        + Nouvelle discussion
      </button>
    </header>

    <div v-if="!isAuthenticated" class="forum-guest-banner">
      <div>
        <span class="forum-chip">Lecture ouverte</span>
        <h2>Vous pouvez consulter tous les forums sans compte.</h2>
        <p>Pour publier une discussion ou repondre, connectez-vous ou creez votre compte UpcycleConnect.</p>
      </div>
      <div class="forum-guest-banner__actions">
        <RouterLink class="btn-secondary forum-guest-button" to="/connexion">Connexion</RouterLink>
        <RouterLink class="btn-main-action forum-guest-button" to="/inscription">Creer un compte</RouterLink>
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
              <span>{{ forum.topics.length }}</span>
            </div>
            <p>{{ forum.description }}</p>
            <small>{{ latestTopicLabel(forum) }}</small>
          </button>
        </div>
      </aside>

      <div class="forum-content">
        <article class="section-container forum-main" v-if="selectedForum">
          <div class="forum-main__head">
            <div>
              <span class="forum-chip">{{ selectedForum.name }}</span>
              <h2>{{ selectedForum.description }}</h2>
            </div>
            <button class="btn-secondary" type="button" @click="openCreateTopic">
              Ecrire un sujet
            </button>
          </div>

          <div v-if="showComposer" class="forum-composer">
            <div class="forum-composer__head">
              <strong>{{ selectedTopic ? "Repondre a la discussion" : "Nouvelle discussion" }}</strong>
              <button class="forum-link" type="button" @click="closeComposer">Annuler</button>
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
              <span class="forum-helper">Publie sous le nom {{ userName }}</span>
              <button class="btn-main-action" type="button" @click="submitPost">
                {{ selectedTopic ? "Envoyer la reponse" : "Publier la discussion" }}
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
              <button type="button" class="forum-topic__summary" @click="toggleTopic(topic.id)">
                <div>
                  <span class="forum-chip forum-chip--muted">{{ topic.tag }}</span>
                  <h3>{{ topic.title }}</h3>
                  <p>{{ topic.preview }}</p>
                </div>
                <div class="forum-topic__meta">
                  <strong>{{ topic.messages.length }} msg</strong>
                  <small>{{ topic.lastActivity }}</small>
                </div>
              </button>

              <div v-if="topic.id === selectedTopicId" class="forum-thread">
                <div
                  v-for="message in topic.messages"
                  :key="message.id"
                  class="forum-message"
                  :class="{ 'is-me': message.author === userName }"
                >
                  <div class="forum-message__avatar">
                    {{ initialsFor(message.author) }}
                  </div>
                  <div class="forum-message__body">
                    <div class="forum-message__meta">
                      <strong>{{ message.author }}</strong>
                      <span>{{ message.role }}</span>
                      <small>{{ message.postedAt }}</small>
                    </div>
                    <p>{{ message.content }}</p>
                  </div>
                </div>

                <div class="forum-thread__footer">
                  <button class="btn-view" type="button" @click="replyToTopic(topic.id)">
                    Repondre
                  </button>
                </div>
              </div>
            </article>
          </div>

          <div v-if="filteredTopics.length === 0" class="forum-empty">
            <p>Aucune discussion ne correspond a votre recherche.</p>
            <button class="btn-secondary" type="button" @click="openCreateTopic">Creer un sujet</button>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import { RouterLink, useRouter } from "vue-router";

const router = useRouter();
const userId = localStorage.getItem("userId") || "guest";
const userPrenom = localStorage.getItem("userPrenom") || "";
const userNom = localStorage.getItem("userNom") || "";
const userName = `${userPrenom} ${userNom}`.trim() || "Utilisateur";
const storageKey = `upcycleconnect-forums-${userId}`;
const isAuthenticated = computed(() => Boolean(localStorage.getItem("userToken") && localStorage.getItem("userId")));

function buildDefaultForums() {
  return [
    {
      id: "general",
      name: "Entraide generale",
      description: "Questions de demarrage, usages de la plateforme et entraide entre membres.",
      topics: [
        {
          id: "general-bienvenue",
          title: "Comment bien preparer sa premiere annonce ?",
          tag: "Debutant",
          preview: "Photos, description, etat de l'objet: les points a verifier avant publication.",
          lastActivity: "Aujourd'hui",
          messages: [
            {
              id: "general-bienvenue-1",
              author: "Equipe UpcycleConnect",
              role: "Moderation",
              postedAt: "Aujourd'hui - 09:10",
              content:
                "Pensez a preciser l'etat de l'objet, sa matiere et sa disponibilite. Une annonce claire est validee plus vite."
            },
            {
              id: "general-bienvenue-2",
              author: "Marie Lambert",
              role: "Particulier",
              postedAt: "Aujourd'hui - 09:42",
              content:
                "Je confirme, ajouter 2 ou 3 photos nettes m'a aidee a recevoir des retours beaucoup plus rapidement."
            }
          ]
        }
      ]
    },
    {
      id: "materiaux",
      name: "Materiaux & recuperation",
      description: "Trouver les bons flux, partager les pieces utiles et identifier les materiaux reutilisables.",
      topics: [
        {
          id: "materiaux-bois",
          title: "Quels objets en bois partent le plus vite ?",
          tag: "Bois",
          preview: "Retours d'experience sur les lots de meubles, chutes et planches a remettre en circulation.",
          lastActivity: "Hier",
          messages: [
            {
              id: "materiaux-bois-1",
              author: "Atelier Renouveau",
              role: "Artisan",
              postedAt: "Hier - 18:20",
              content:
                "Les petites commodes, tabourets et plateaux de table trouvent vite preneur quand les dimensions sont indiquees."
            }
          ]
        }
      ]
    },
    {
      id: "projets",
      name: "Projets & transformations",
      description: "Montrer ses creations, demander des avis et partager les avant/apres.",
      topics: [
        {
          id: "projets-jean",
          title: "Idees pour transformer un vieux jean en accessoire",
          tag: "Textile",
          preview: "Besoin d'inspirations simples a realiser a la maison ou en atelier.",
          lastActivity: "Cette semaine",
          messages: [
            {
              id: "projets-jean-1",
              author: "Lea Dupont",
              role: "Particulier",
              postedAt: "Cette semaine - 14:05",
              content:
                "Le sac cabas marche tres bien, surtout si vous gardez les poches avant pour le cote pratique."
            }
          ]
        }
      ]
    }
  ];
}

function loadForums() {
  const raw = localStorage.getItem(storageKey);
  if (!raw) return buildDefaultForums();

  try {
    const parsed = JSON.parse(raw);
    return Array.isArray(parsed) && parsed.length ? parsed : buildDefaultForums();
  } catch {
    return buildDefaultForums();
  }
}

const forums = ref(loadForums());
const search = ref("");
const selectedForumId = ref(forums.value[0]?.id || "");
const selectedTopicId = ref(forums.value[0]?.topics[0]?.id || "");
const showComposer = ref(false);
const draftTopicTitle = ref("");
const draftMessage = ref("");

const selectedForum = computed(() => forums.value.find((forum) => forum.id === selectedForumId.value) || null);
const selectedTopic = computed(() => selectedForum.value?.topics.find((topic) => topic.id === selectedTopicId.value) || null);

const filteredForums = computed(() => {
  const term = search.value.toLowerCase();
  if (!term) return forums.value;

  return forums.value.filter((forum) => {
    const forumMatch = forum.name.toLowerCase().includes(term) || forum.description.toLowerCase().includes(term);
    const topicMatch = forum.topics.some(
      (topic) => topic.title.toLowerCase().includes(term) || topic.preview.toLowerCase().includes(term)
    );
    return forumMatch || topicMatch;
  });
});

const filteredTopics = computed(() => {
  if (!selectedForum.value) return [];
  const term = search.value.toLowerCase();
  if (!term) return selectedForum.value.topics;

  return selectedForum.value.topics.filter(
    (topic) =>
      topic.title.toLowerCase().includes(term) ||
      topic.preview.toLowerCase().includes(term) ||
      topic.messages.some((message) => message.content.toLowerCase().includes(term))
  );
});

const totalTopics = computed(() => forums.value.reduce((count, forum) => count + forum.topics.length, 0));
const totalMessages = computed(() =>
  forums.value.reduce(
    (count, forum) => count + forum.topics.reduce((topicCount, topic) => topicCount + topic.messages.length, 0),
    0
  )
);

function persistForums() {
  localStorage.setItem(storageKey, JSON.stringify(forums.value));
}

function latestTopicLabel(forum) {
  if (!forum.topics.length) return "Aucune discussion";
  return forum.topics[0].title;
}

function initialsFor(name) {
  return name
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() || "")
    .join("");
}

function requireAuth() {
  if (isAuthenticated.value) return true;
  router.push({ path: "/connexion", query: { redirect: "/forums" } });
  return false;
}

function selectForum(forumId) {
  selectedForumId.value = forumId;
  selectedTopicId.value = forums.value.find((forum) => forum.id === forumId)?.topics[0]?.id || "";
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

function submitPost() {
  if (!requireAuth() || !selectedForum.value || !draftMessage.value.trim()) return;

  const now = new Date();
  const formattedDate = now.toLocaleDateString("fr-FR", { day: "numeric", month: "short" });
  const formattedTime = now.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" });

  if (selectedTopic.value) {
    selectedTopic.value.messages.push({
      id: `msg-${Date.now()}`,
      author: userName,
      role: "Particulier",
      postedAt: `${formattedDate} - ${formattedTime}`,
      content: draftMessage.value.trim()
    });
    selectedTopic.value.lastActivity = "A l'instant";
  } else {
    if (draftTopicTitle.value.trim().length < 4) return;

    const newTopic = {
      id: `topic-${Date.now()}`,
      title: draftTopicTitle.value.trim(),
      tag: "Discussion",
      preview: draftMessage.value.trim().slice(0, 110),
      lastActivity: "A l'instant",
      messages: [
        {
          id: `msg-${Date.now()}-root`,
          author: userName,
          role: "Particulier",
          postedAt: `${formattedDate} - ${formattedTime}`,
          content: draftMessage.value.trim()
        }
      ]
    };

    selectedForum.value.topics.unshift(newTopic);
    selectedTopicId.value = newTopic.id;
  }

  persistForums();
  closeComposer();
}

watch(filteredForums, (nextForums) => {
  if (!nextForums.length) return;
  if (!nextForums.some((forum) => forum.id === selectedForumId.value)) {
    selectedForumId.value = nextForums[0].id;
    selectedTopicId.value = nextForums[0].topics[0]?.id || "";
  }
});
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
</style>
