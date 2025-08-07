<template>
  <div class="conversation-info-sidebar">
    <div class="conversation-info-sidebar__header">
      <div class="conversation-info-sidebar__header__back-and-title">
        <button v-if="selectedMember" class="back-btn" @click="backToGroupInfo">
          <svg viewBox="0 0 24 24" fill="none" class="back-btn-icon">
            <path
              fill-rule="evenodd"
              clip-rule="evenodd"
              d="M11.7071 4.29289C12.0976 4.68342 12.0976 5.31658 11.7071 5.70711L6.41421 11H20C20.5523 11 21 11.4477 21 12C21 12.5523 20.5523 13 20 13H6.41421L11.7071 18.2929C12.0976 18.6834 12.0976 19.3166 11.7071 19.7071C11.3166 20.0976 10.6834 20.0976 10.2929 19.7071L3.29289 12.7071C3.10536 12.5196 3 12.2652 3 12C3 11.7348 3.10536 11.4804 3.29289 11.2929L10.2929 4.29289C10.6834 3.90237 11.3166 3.90237 11.7071 4.29289Z"
              fill="var(--color-secondary)"
            />
          </svg>
        </button>
        <span class="text-body">
          {{
            selectedMember
              ? "Contact info"
              : conversation.type === "private"
                ? "Contact info"
                : "Group info"
          }}
        </span>
      </div>
      <button aria-label="Close" class="close-btn" @click="emit('close')">
        <svg viewBox="0 0 24 24" fill="none" class="close-btn-icon">
          <path
            d="M20.7457 3.32851C20.3552 2.93798 19.722 2.93798 19.3315 3.32851L12.0371 10.6229L4.74275 3.32851C4.35223 2.93798 3.71906 2.93798 3.32854 3.32851C2.93801 3.71903 2.93801 4.3522 3.32854 4.74272L10.6229 12.0371L3.32856 19.3314C2.93803 19.722 2.93803 20.3551 3.32856 20.7457C3.71908 21.1362 4.35225 21.1362 4.74277 20.7457L12.0371 13.4513L19.3315 20.7457C19.722 21.1362 20.3552 21.1362 20.7457 20.7457C21.1362 20.3551 21.1362 19.722 20.7457 19.3315L13.4513 12.0371L20.7457 4.74272C21.1362 4.3522 21.1362 3.71903 20.7457 3.32851Z"
            fill="var(--color-tertiary)"
          />
        </svg>
      </button>
    </div>
    <template v-if="selectedMember">
      <div class="conversation-info-sidebar__content">
        <img
          :src="
            selectedMember?.profilePicture
              ? backendBaseUrl + selectedMember.profilePicture
              : defaultProfilePicture
          "
          class="profile__picture"
          alt="Profile picture"
        >
        <div class="profile__username-field">
          <input
            class="profile__username conversation-info__username"
            :value="
              selectedMember?.userId === user.userId
                ? 'You'
                : selectedMember?.username
            "
            disabled
          >
        </div>
      </div>
    </template>
    <template v-else-if="conversation.type === 'group'">
      <div class="conversation-info-sidebar__content">
        <label class="profile__picture-field">
          <img
            :src="editedPhotoUrl"
            alt="Group picture"
            class="profile__picture"
            :class="{ 'profile__picture--error': photoError }"
          >
          <svg viewBox="0 0 24 24" fill="none" class="profile__picture-icon">
            <path
              d="M20.1498 7.93997L8.27978 19.81C7.21978 20.88 4.04977 21.3699 3.32977 20.6599C2.60977 19.9499 3.11978 16.78 4.17978 15.71L16.0498 3.84C16.5979 3.31801 17.3283 3.03097 18.0851 3.04019C18.842 3.04942 19.5652 3.35418 20.1004 3.88938C20.6356 4.42457 20.9403 5.14781 20.9496 5.90463C20.9588 6.66146 20.6718 7.39189 20.1498 7.93997V7.93997Z"
              stroke="var(--color-secondary)"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
          <input
            type="file"
            accept="image/*"
            style="display: none"
            @change="onPhotoChange"
          >
        </label>
        <div class="profile__username-field">
          <div class="conversation-info__group-name-input">
            <input
              ref="nameInput"
              v-model="editedName"
              class="profile__username conversation-info__group-name"
              :class="{
                'conversation-info__group-name--error': nameError,
              }"
              :disabled="!editing"
              @keydown.enter="saveName"
            >
            <button
              class="profile__username-icon-button"
              @click="editing ? saveName() : startEditing()"
            >
              <svg
                v-if="!editing"
                viewBox="0 0 24 24"
                fill="none"
                class="profile__username-icon"
              >
                <path
                  d="M20.1498 7.93997L8.27978 19.81C7.21978 20.88 4.04977 21.3699 3.32977 20.6599C2.60977 19.9499 3.11978 16.78 4.17978 15.71L16.0498 3.84C16.5979 3.31801 17.3283 3.03097 18.0851 3.04019C18.842 3.04942 19.5652 3.35418 20.1004 3.88938C20.6356 4.42457 20.9403 5.14781 20.9496 5.90463C20.9588 6.66146 20.6718 7.39189 20.1498 7.93997V7.93997Z"
                  stroke="var(--color-tertiary)"
                  stroke-width="1.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
              <svg
                v-else
                viewBox="0 0 24 24"
                fill="none"
                class="profile__username-icon"
              >
                <path
                  d="M4 12.6111L8.92308 17.5L20 6.5"
                  :stroke="
                    nameError ? 'var(--color-error)' : 'var(--color-tertiary)'
                  "
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </button>
          </div>
        </div>
        <div class="group-members">
          <span class="text-body" style="color: var(--color-tertiary)">{{ conversation.members.length }} Members</span>
          <button
            class="add-member-btn"
            @click="$emit('add-member-modal-open')"
          >
            <div class="add-member-btn-icon">
              <svg width="24" height="24" viewBox="0 0 1920 1920" fill="none">
                <path
                  d="M866.332 213v653.332H213v186.666h653.332v653.332h186.666v-653.332h653.332V866.332h-653.332V213z"
                  fill-rule="evenodd"
                  fill="var(--color-secondary)"
                />
              </svg>
            </div>
            <span class="text-body">Add member</span>
          </button>
          <ul class="users">
            <li
              v-for="member in sortedMembers"
              :key="member.userId"
              class="user"
            >
              <button class="user__button" @click="showContactInfo(member)">
                <img
                  :src="
                    member?.profilePicture
                      ? backendBaseUrl + member.profilePicture
                      : defaultProfilePicture
                  "
                  alt="Profile picture"
                  class="user__picture"
                >
                <span class="text-body">{{
                  member.userId === user.userId ? "You" : member.username
                }}</span>
              </button>
            </li>
          </ul>
          <button class="leave-group-btn" @click="leaveGroup">
            <div class="leave-group-btn-icon">
              <svg width="24" height="24" viewBox="0 0 20 20" fill="none">
                <path
                  fill-rule="evenodd"
                  clip-rule="evenodd"
                  d="M15.6666 8L17.75 10.5L15.6666 8Z"
                  stroke="var(--color-error)"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  fill-rule="evenodd"
                  clip-rule="evenodd"
                  d="M15.6666 13L17.75 10.5L15.6666 13Z"
                  stroke="var(--color-error)"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M16.5 10.5L10 10.5"
                  stroke="var(--color-error)"
                  stroke-width="2"
                  stroke-linecap="round"
                />
                <line
                  x1="4"
                  y1="3.5"
                  x2="13"
                  y2="3.5"
                  stroke="var(--color-error)"
                  stroke-width="2"
                  stroke-linecap="round"
                />
                <line
                  x1="4"
                  y1="17.5"
                  x2="13"
                  y2="17.5"
                  stroke="var(--color-error)"
                  stroke-width="2"
                  stroke-linecap="round"
                />
                <path
                  d="M13 3.5V7.5"
                  stroke="var(--color-error)"
                  stroke-width="2"
                  stroke-linecap="round"
                />
                <path
                  d="M13 13.5V17.5"
                  stroke="var(--color-error)"
                  stroke-width="2"
                  stroke-linecap="round"
                />
                <path
                  d="M4 3.5L4 17.5"
                  stroke="var(--color-error)"
                  stroke-width="2"
                  stroke-linecap="round"
                />
              </svg>
            </div>
            <span class="text-body">Leave group</span>
          </button>
        </div>
      </div>
    </template>
    <template v-else>
      <div class="conversation-info-sidebar__content">
        <img
          :src="
            otherUser?.profilePicture
              ? backendBaseUrl + otherUser.profilePicture
              : defaultProfilePicture
          "
          class="profile__picture"
          alt="Profile picture"
        >
        <div class="profile__username-field">
          <input
            class="profile__username conversation-info__username"
            :value="otherUser?.username"
            disabled
          >
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, nextTick, watch } from "vue";
import { backendBaseUrl } from "@/services/baseUrl";
import defaultProfilePicture from "@/assets/default-profile-picture.jpg";
import defaultGroupPicture from "@/assets/default-group-picture.jpg";
import api from "@/services/api";

const props = defineProps({
  conversation: Object,
  user: Object,
});

const emit = defineEmits([
  "close",
  "group-name-updated",
  "group-updated",
  "left-group",
  "add-member-modal-open",
]);

const otherUser = computed(() =>
  props.conversation.type === "private"
    ? props.conversation.participants.find(
        (u) => u.userId !== props.user.userId,
      )
    : null,
);

const sortedMembers = computed(() => {
  if (!props.conversation.members) return [];

  return [
    ...props.conversation.members.filter((m) => m.userId === props.user.userId),
    ...props.conversation.members.filter((m) => m.userId !== props.user.userId),
  ];
});

const editing = ref(false);
const editedName = ref(props.conversation.name);
const nameInput = ref(null);
const nameError = ref(false);

const selectedMember = ref(null);

function startEditing() {
  editing.value = true;
  nextTick(() => nameInput.value?.focus());
}

async function saveName() {
  if (editedName.value === props.conversation.name) {
    editing.value = false;
    return;
  }

  nameError.value = false;

  try {
    await api.put(`/groups/${props.conversation.conversationId}/name`, {
      name: editedName.value,
    });

    emit("group-name-updated", editedName.value);

    editing.value = false;
  } catch (e) {
    nameError.value = true;
  }
}

watch(editedName, () => {
  nameError.value = false;
});

const photoError = ref(false);
const editedPhotoUrl = computed(() =>
  props.conversation.photo
    ? backendBaseUrl + props.conversation.photo
    : defaultGroupPicture,
);

async function onPhotoChange(event) {
  const file = event.target.files[0];
  if (!file) return;

  const formData = new FormData();
  formData.append("image", file);

  photoError.value = false;

  try {
    const response = await api.put(
      `/groups/${props.conversation.conversationId}/photo`,
      formData,
      {
        headers: { "Content-Type": "multipart/form-data" },
      },
    );
    emit("group-updated", response.data);
  } catch (e) {
    photoError.value = true;
  }
}

async function leaveGroup() {
  try {
    await api.delete(`/groups/${props.conversation.conversationId}/members/me`);

    emit("close");
    emit("left-group");
  } catch (e) {
    console.error(e);
  }
}

function showContactInfo(member) {
  selectedMember.value = member;
}

function backToGroupInfo() {
  selectedMember.value = null;
}
</script>

<style scoped>
.conversation-info-sidebar {
  position: fixed;
  right: 0;
  top: 0;
  width: 30vw;
  height: 100vh;
  background: var(--color-background);
  z-index: 10001;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  border-left: 1px solid var(--color-quaternary);
}

.conversation-info-sidebar__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem;
}

.conversation-info-sidebar__header__back-and-title {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
}

.back-btn,
.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 50%;
  padding: 0.3rem;
  background-color: inherit;
}

.back-btn:hover,
.close-btn:hover {
  background-color: var(--color-quaternary);
  transition: background 0.1s;
}

.back-btn-icon,
.close-btn-icon {
  width: 20px;
  height: 20px;
}

.conversation-info-sidebar__content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1rem;
  gap: 2.5rem;
}

.conversation-info__group-name-input {
  display: flex;
}

.conversation-info__username,
.conversation-info__group-name {
  text-align: center;
  font-weight: 600;
  font-size: 1.25rem;
}

.conversation-info__group-name--error {
  color: var(--color-error);
}

.group-members {
  width: 100%;
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: 0.25rem;
}

.group-members > .text-body {
  align-self: flex-start;
  padding: 1rem;
}

.add-member-btn,
.leave-group-btn {
  display: flex;
  align-items: center;
  gap: 1rem;
  border: none;
  border-radius: 12px;
  width: 100%;
  padding: 0.5rem 0.75rem;
  background-color: inherit;
  color: var(--color-secondary);
  transition: background-color 0.1s;
}

.add-member-btn:hover,
.add-member-btn:focus,
.leave-group-btn:hover,
.leave-group-btn:focus {
  outline: none;
  background-color: var(--color-quaternary);
}

.add-member-btn-icon,
.leave-group-btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border: 2px transparent;
  border-radius: 50%;
  object-fit: cover;
  background-color: var(--color-primary);
}

.leave-group-btn-icon {
  background: none;
}
</style>
