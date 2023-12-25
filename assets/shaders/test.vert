#version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;
layout (location = 2) in vec3 aNormal;

out vec3 ModelPos;

out vec2 TexCoord;

out vec3 Normal;
out vec3 FragPos;

uniform mat4 model;
uniform mat4 view;
uniform mat4 proj;

void main() {
	FragPos = vec3(model*vec4(aPos,1.0));

	gl_Position = proj*view*model*vec4(FragPos,1.0f);
	TexCoord = vec2(aTexCoord.x, 1.0f - aTexCoord.y);
	ModelPos = vec3(model[3]);

	Normal = aNormal;
}
