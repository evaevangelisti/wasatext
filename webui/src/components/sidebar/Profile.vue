<template>
    <div class="profile">
        <label class="profile__picture-field">
            <img
                :src="
                    user?.profilePicture
                        ? backendBaseUrl + user.profilePicture
                        : defaultProfilePicture
                "
                alt="Profile picture"
                class="profile__picture"
                :class="{ 'profile__picture--error': photoError }"
            />
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
                @change="onPhotoChange"
                style="display: none"
            />
        </label>
        <div class="profile__username-field">
            <span class="text-secondary">Username</span>
            <div
                class="profile__username-input"
                :class="{ 'profile__username-input--error': usernameError }"
            >
                <input
                    id="profile__username"
                    ref="usernameInput"
                    v-model="editedUsername"
                    :disabled="!editing"
                />
                <button
                    ref="usernameButton"
                    class="profile__username-icon-button"
                    @click="editing ? saveUsername() : startEditing()"
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
                            stroke="var(--color-secondary)"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        />
                    </svg>
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, watch, nextTick } from "vue";
import api from "@/services/api";
import { backendBaseUrl } from "@/services/baseUrl";

import defaultProfilePicture from "@/assets/default-profile-picture.jpg";

const props = defineProps({
    user: Object,
});

const emit = defineEmits(["profile-updated"]);

const editing = ref(false);
const editedUsername = ref(props.user?.username || "");
const usernameInput = ref(null);
const usernameButton = ref(null);
const usernameError = ref(false);

watch(
    () => props.user?.username,
    (newValue) => {
        editedUsername.value = newValue || "";
        editing.value = false;
    },
);

function startEditing() {
    editing.value = true;
    nextTick(() => {
        usernameInput.value?.focus();
    });
}

async function saveUsername() {
    if (editedUsername.value === props.user?.username) {
        usernameInput.value?.blur();
        usernameButton.value?.blur();
        editing.value = false;
        return;
    }

    usernameError.value = false;

    try {
        const response = await api.put("/me/username", {
            username: editedUsername.value,
        });

        emit("profile-updated", response.data);
        editing.value = false;

        usernameInput.value?.blur();
        usernameButton.value?.blur();
    } catch (e) {
        console.log(e);
        usernameError.value = true;

        usernameInput.value?.blur();
        usernameButton.value?.blur();
    }
}

const photoError = ref(false);

async function onPhotoChange(event) {
    const file = event.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append("image", file);

    photoError.value = false;

    try {
        const response = await api.put("/me/photo", formData, {
            headers: { "Content-Type": "multipart/form-data" },
        });

        emit("profile-updated", response.data);
    } catch (e) {
        photoError.value = true;
    }
}
</script>

<style setup>
.profile {
    display: flex;
    align-items: center;
    flex-direction: column;
    gap: 2.5rem;
    padding: 3rem;
}

.profile__picture-field {
    position: relative;
    display: inline-block;
    cursor: pointer;
}

.profile__picture {
    border: 2px solid transparent;
    border-radius: 50%;
    width: 128px;
    height: 128px;
    object-fit: cover;
    transition: opacity 0.1s;
}

.profile__picture--error {
    border-color: var(--color-error) !important;
}

.profile__picture-field:hover .profile__picture {
    opacity: 0.5;
}

.profile__picture-icon {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    pointer-events: none;
    width: 24px;
    height: 24px;
    opacity: 0;
    transition: opacity 0.1s;
}

.profile__picture-field:hover .profile__picture-icon {
    opacity: 1;
}

.profile__username-field {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    width: 100%;
}

.profile__username-input {
    display: flex;
    border-bottom: 2px solid var(--color-tertiary);
}

.profile__username-input--error {
    border-bottom-color: var(--color-error);
}

.profile__username-input:focus-within {
    border-bottom: 2px solid var(--color-primary);
}

#profile__username {
    border: none;
    width: 100%;
    padding: 0.5rem 0rem;
    background-color: inherit;
    font-size: 1rem;
    color: var(--color-secondary);
}

.profile__username-icon-button {
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    padding: 0.2rem;
    background-color: inherit;
}

.profile__username-icon {
    width: 24px;
    height: 24px;
}
</style>
