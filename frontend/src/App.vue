<script setup>
import {ref, onMounted} from 'vue'
import {userService} from './services/api'
import UserFormModal from './components/UserFormModal.vue'

const users = ref([])
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const isModalOpen = ref(false)
const selectedUser = ref(null)

async function fetchUsers(){
  loading.value = true
  errorMessage.value = ''
  try{
    users.value = await userService.getAll()
  } catch (err){
    errorMessage.value = err.message
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  selectedUser.value = null
  successMessage.value = ''
  isModalOpen.value = true
}

function openEditModal(user) {
  selectedUser.value = user
  successMessage.value = ''
  isModalOpen.value = true
}

function handleSaved(message) {
  isModalOpen.value = false
  successMessage.value = message
  fetchUsers()
}

async function handleDelete(user) {
  if (!confirm(`Deseja realmente excluir o usuário "${user.name}"?`)) return

  errorMessage.value = ''
  successMessage.value = ''
  try {
    await userService.remove(user.id)
    successMessage.value = 'Usuário excluído com sucesso!'
    fetchUsers()
  } catch (err) {
    errorMessage.value = err.message
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <UserFormModal
      v-if="isModalOpen"
      :user-to-edit="selectedUser"
      @saved="handleSaved"
      @close="isModalOpen = false"
    />
  <main>
    <h1>User Management</h1>
    <button @click="openCreateModal">+ Novo Usuário</button>

    <p v-if="successMessage" class="success">{{ successMessage }}</p>

    <p v-if="loading">Carregando usuários...</p>
    <p v-if="errorMessage">{{ errorMessage }}</p>

    <table v-if="!loading && users.length > 0" border="1">
      <thead>
        <tr>
          <th>Nome</th>
          <th>E-mail</th>
          <th>Login</th>
          <th>Status</th>
          <th>Ações</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in users" :key="user.id">
          <td>{{ user.name }}</td>
          <td>{{ user.email }}</td>
          <td>{{ user.login }}</td>
          <td>{{ user.active ? 'Ativo' : 'Inativo' }}</td>
          <td>
            <button @click="openEditModal(user)">Editar</button>
            <button @click="handleDelete(user)">Excluir</button>
          </td>
        </tr>
      </tbody>
    </table>

    <p v-if="!loading && users.length === 0">Nenhum usuário cadastrado.</p>
  </main>
</template>

<style scoped>
header {
  line-height: 1.5;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
}
</style>
