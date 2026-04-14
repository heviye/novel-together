import React from 'react';
import { View, Text, Button, StyleSheet } from 'react-native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { useAuth } from '../context/AuthContext';

type RootStackParamList = {
  Home: undefined;
  Login: undefined;
  Register: undefined;
  NovelList: undefined;
  NovelDetail: { novelId: string };
  WriteChapter: { novelId: string };
  Chapter: { chapterId: string };
  Profile: undefined;
};

type HomeScreenProps = {
  navigation: NativeStackNavigationProp<RootStackParamList, 'Home'>;
};

export default function HomeScreen({ navigation }: HomeScreenProps) {
  const { user, loading } = useAuth();

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Novel Together</Text>
      <Text style={styles.subtitle}>协作创作平台</Text>
      
      {user ? (
        <>
          <Text style={styles.welcome}>欢迎, {user.username}!</Text>
          <View style={styles.buttonContainer}>
            <Button
              title="浏览小说"
              onPress={() => navigation.navigate('NovelList')}
            />
            <Button
              title="我的资料"
              onPress={() => navigation.navigate('Profile')}
            />
          </View>
        </>
      ) : (
        <View style={styles.buttonContainer}>
          <Button
            title="登录"
            onPress={() => navigation.navigate('Login')}
          />
          <Button
            title="注册"
            onPress={() => navigation.navigate('Register')}
          />
          <Button
            title="浏览小说"
            onPress={() => navigation.navigate('NovelList')}
          />
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    marginBottom: 10,
  },
  subtitle: {
    fontSize: 16,
    color: '#666',
    marginBottom: 40,
  },
  welcome: {
    fontSize: 18,
    marginBottom: 20,
  },
  buttonContainer: {
    gap: 15,
    width: '80%',
  },
});