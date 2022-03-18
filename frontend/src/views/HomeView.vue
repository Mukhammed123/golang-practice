<script setup lang="ts">
import axios from 'axios';
import { ref } from 'vue';

const fileNames = ref([]);
const loadFunc = async () => {
  const config = {
    method: 'get',
    url: `http://localhost:8081/people`,
    headers: {
      "Content-Type": "application/json",
    },}
  const response = await axios(config);
  fileNames.value = response.data;
  console.log(response);
}
loadFunc();
const uploadedFile = ref(null)
const createUser = async () => {

  const config = {
    method: 'post',
    url: `http://localhost:8081/people`,
    headers: {
      "Content-Type": "application/json",
    },
    data: {
      Name: "Mikeee",
      Email: "mikeee@gmail.com"
    },
    onUploadProgress: (uploadEvent) => {
      console.log("Uploaded Progress: " + Math.round(uploadEvent.loaded / uploadEvent.total)*100 + '%');
    }
  }
  // const response = await axios(config);
  console.log(response);
}

const uploadFile = async (event) => {
  uploadedFile.value = event.target.files[0];
}

const sendFile = async () => {
  const buffer = await uploadedFile.value.arrayBuffer();
    let byteObj = new Int8Array(buffer);
    let byteArray: number[] = []
    console.log(byteObj)
    for (var key in byteObj) {
      byteArray.push(byteObj[key])
    }
  const config = {
    method: 'post',
    url: `http://localhost:8081/files`,
    headers: {
      "Content-Type": "application/json",
    },
    data: {
      Name: uploadedFile.value.name,
      Type: uploadedFile.value.type,
      Content: byteArray
    },
    onUploadProgress: (uploadEvent) => {
      console.log("Uploaded Progress: " + Math.round(uploadEvent.loaded / uploadEvent.total)*100 + '%');
    }
  }

  axios(config);
}

// const downloadFile = async () => {
  // const config = {
  //   method: 'post',
  //   url: `http://localhost:8081/downloadFile`,
  //   headers: {
  //     "Content-Disposition": "attachment;filename=123.png",
  //     "Content-Type": "application/octet-stream"
  //   }
  // }
  // axios(config)
// }

</script>

<template>
  <main>
    <form :action="`http://localhost:8081/files/${uploadedFile ? uploadedFile.name : ''}`" method="post" enctype="multipart/form-data" >
      <input type="file" name="myFile" @input="uploadFile">
      <input type="submit" value="Uploading file...">
    </form>
    <div v-for="fileName in fileNames" :key="fileName" style="margin: 1em;">
    <a :href="`http://localhost:8081/downloadFile/${fileName}`">{{fileName}}</a >
    </div >
  </main>
</template>
