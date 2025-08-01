<template>
    <div class="dashboard">
        <Sidebar
            :user="user"
            @profile-updated="$emit('profile-updated', $event)"
        />
        <Conversation />
    </div>
</template>

<script setup>
import Sidebar from "@/components/sidebar/Sidebar.vue";
import Conversation from "@/components/conversation/Conversation.vue";

const props = defineProps({
    user: Object,
});

const emit = defineEmits(["profile-updated"]);
</script>

<style setup>
.dashboard {
    display: flex;
    height: 100vh;
}
</style>
